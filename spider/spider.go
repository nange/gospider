package spider

import (
	"context"
	"database/sql"

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
			ctx = newContext(ctxCtl, task, collectors[i], collectors[i+1])
		} else {
			ctx = newContext(ctxCtl, task, collectors[i], nil)
		}
		if task.OutputConfig.Type == common.OutputTypeMySQL {
			ctx.setOutputDB(db)
		}

		addCallback(ctx, task.Rule.Nodes[i])
	}

	headCtx := newContext(ctxCtl, task, c, collectors[0])
	if err := task.Rule.Head(headCtx); err != nil {
		logrus.Errorf("exec rule head func err:%#v", err)
		return errors.WithStack(err)
	}

	go func() {
		defer db.Close()
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

func addCallback(ctx *Context, node *Node) {
	if node.OnRequest != nil {
		ctx.c.OnRequest(func(req *colly.Request) {
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
