package spider

import (
	"context"
	"database/sql"
	"fmt"
	"runtime/debug"

	"github.com/gocolly/colly"
	"github.com/nange/gospider/common"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func Run(task *Task, retCh chan<- common.MTS) error {
	var db *sql.DB
	var err error
	if task.OutputConfig.Type == common.OutputTypeMySQL {
		db, err = newDB(task.TaskConfig.OutputConfig.MySQLConf)
		if err != nil {
			logrus.Errorf("newDB failed! err:%#v", err)
			return err
		}
	}

	c, err := newCollector(task.TaskConfig)
	if err != nil {
		logrus.Errorf("new collector err:%+v", err)
		return err
	}

	nodesLen := len(task.Rule.Nodes)
	collectors := make([]*colly.Collector, 0, nodesLen)
	for i := 0; i < len(task.Rule.Nodes); i++ {
		nextC := c.Clone()
		collectors = append(collectors, nextC)
	}

	ctxCtl, cancel := context.WithCancel(context.Background())
	if err := addTaskCtrl(task.ID, cancel); err != nil {
		return errors.Wrapf(err, "addTaskCtrl failed")
	}

	for i := 0; i < nodesLen; i++ {
		var ctx *Context
		if i != nodesLen-1 {
			ctx = newContext(ctxCtl, cancel, task, collectors[i], collectors[i+1])
		} else {
			ctx = newContext(ctxCtl, cancel, task, collectors[i], nil)
		}
		if task.OutputConfig.Type == common.OutputTypeMySQL {
			ctx.setOutputDB(db)
		}

		addCallback(ctx, task.Rule.Nodes[i])
	}

	headCtx := newContext(ctxCtl, cancel, task, c, collectors[0])
	headWrapper := func(ctx *Context) (err error) {
		defer func() {
			if e := recover(); e != nil {
				if v, ok := e.(error); ok {
					err = v
				} else {
					str := fmt.Sprintf("%v", e)
					err = errors.New(str)
				}
				logrus.Errorf("Head unexcepted exited, err: %+v, stack:\n%s", e, string(debug.Stack()))
			}
		}()
		return task.Rule.Head(ctx)
	}
	if err := headWrapper(headCtx); err != nil {
		logrus.Errorf("exec rule head func err:%#v", err)
		return errors.WithStack(err)
	}

	go func() {
		if db != nil {
			defer db.Close()
		}

		for i := 0; i < nodesLen; i++ {
			collectors[i].Wait()
			logrus.Infof("task:%s %d step completed...", task.Name, i+1)
		}

		if err := CancelTask(task.ID); err == nil {
			retCh <- common.MTS{ID: task.ID, Status: common.TaskStatusCompleted}
		}
		logrus.Infof("task:%s run completed...", task.Name)
	}()

	return nil
}

func cbDefer(ctx *Context, info string) {
	if e := recover(); e != nil {
		logrus.Error(info + fmt.Sprintf(", err: %+v, stack:\n%s", e, string(debug.Stack())))
		ctx.ctlCancel()
	}
}

func addCallback(ctx *Context, node *Node) {
	if node.OnRequest != nil {
		ctx.c.OnRequest(func(req *colly.Request) {
			defer cbDefer(ctx, fmt.Sprintf("OnRequest unexcepted exited, url:%s", req.URL.String()))

			newCtx := ctx.cloneWithReq(req)
			select {
			case <-newCtx.ctlCtx.Done():
				logrus.Warnf("request has been canceled in OnRequest, url:%s", newCtx.GetRequest().URL.String())
				newCtx.Abort()
				return
			default:
			}

			node.OnRequest(newCtx, newRequest(req, newCtx))
		})
	}

	if node.OnError != nil {
		ctx.c.OnError(func(res *colly.Response, e error) {
			defer cbDefer(ctx, fmt.Sprintf("OnError unexcepted exited, url:%s", res.Request.URL.String()))

			newCtx := ctx.cloneWithReq(res.Request)
			select {
			case <-newCtx.ctlCtx.Done():
				logrus.Warnf("request has been canceled in OnError, url:%s", newCtx.GetRequest().URL.String())
				return
			default:
			}

			err := node.OnError(newCtx, newResponse(res, newCtx), e)
			if err != nil {
				logrus.Warnf("node.OnError return err:%+v, request url:%s", err, res.Request.URL.String())
			}
		})
	}

	if node.OnResponse != nil {
		ctx.c.OnResponse(func(res *colly.Response) {
			defer cbDefer(ctx, fmt.Sprintf("OnResponse unexcepted exited, url:%s", res.Request.URL.String()))

			newCtx := ctx.cloneWithReq(res.Request)
			select {
			case <-newCtx.ctlCtx.Done():
				logrus.Warnf("request has been canceled in OnResponse, url:%s", newCtx.GetRequest().URL.String())
				return
			default:
			}

			err := node.OnResponse(newCtx, newResponse(res, newCtx))
			if err != nil {
				logrus.Warnf("node.OnResponse return err:%+v, request url:%s", err, res.Request.URL.String())
			}
		})
	}

	if node.OnHTML != nil {
		for selector, fn := range node.OnHTML {
			f := fn
			ctx.c.OnHTML(selector, func(el *colly.HTMLElement) {
				defer cbDefer(ctx, fmt.Sprintf("OnHTML unexcepted exited, selector:%s, url:%s", selector, el.Request.URL.String()))

				newCtx := ctx.cloneWithReq(el.Request)
				select {
				case <-newCtx.ctlCtx.Done():
					logrus.Warnf("request has been canceled in OnHTML, url:%s", newCtx.GetRequest().URL.String())
					return
				default:
				}

				err := f(newCtx, newHTMLElement(el, newCtx))
				if err != nil {
					logrus.Warnf("node.OnHTML:%s return err:%+v, request url:%s", selector, err, el.Request.URL.String())
				}
			})
		}
	}

	if node.OnXML != nil {
		for selector, fn := range node.OnXML {
			f := fn
			ctx.c.OnXML(selector, func(el *colly.XMLElement) {
				defer cbDefer(ctx, fmt.Sprintf("OnXML unexcepted exited, selector:%s, url:%s", selector, el.Request.URL.String()))

				newCtx := ctx.cloneWithReq(el.Request)
				select {
				case <-newCtx.ctlCtx.Done():
					logrus.Warnf("request has been canceled in OnXML, url:%s", newCtx.GetRequest().URL.String())
					return
				default:
				}

				err := f(newCtx, newXMLElement(el, newCtx))
				if err != nil {
					logrus.Warnf("node.OnXML:%s return err:%+v, request url:%s", selector, err, el.Request.URL.String())
				}
			})
		}
	}

	if node.OnScraped != nil {
		ctx.c.OnScraped(func(res *colly.Response) {
			defer cbDefer(ctx, fmt.Sprintf("OnScraped unexcepted exited, url:%s", res.Request.URL.String()))

			newCtx := ctx.cloneWithReq(res.Request)
			select {
			case <-newCtx.ctlCtx.Done():
				logrus.Warnf("request has been canceled in OnScraped, url:%s", newCtx.GetRequest().URL.String())
				return
			default:
			}

			err := node.OnScraped(newCtx, newResponse(res, newCtx))
			if err != nil {
				logrus.Warnf("node.OnScraped return err:%+v, request url:%s", err, res.Request.URL.String())
			}
		})
	}

}
