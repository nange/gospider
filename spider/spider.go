package spider

import (
	"context"
	"database/sql"
	"fmt"
	"runtime/debug"

	"github.com/gocolly/colly"
	"github.com/nange/gospider/common"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// Spider the spider define
type Spider struct {
	task  *Task
	retCh chan<- common.MTS
	db    *sql.DB
}

// New create a new spider object
func New(task *Task, retCh chan<- common.MTS) *Spider {
	return &Spider{
		task:  task,
		retCh: retCh,
	}
}

// SetDB set the underlayer output db
func (s *Spider) SetDB(db *sql.DB) {
	s.db = db
}

// Run run a spider task
func (s *Spider) Run() error {
	c, err := newCollector(s.task.TaskConfig)
	if err != nil {
		log.Errorf("new collector err:%+v", err)
		return err
	}

	nodesLen := len(s.task.Rule.Nodes)
	collectors := make([]*colly.Collector, 0, nodesLen)
	for i := 0; i < len(s.task.Rule.Nodes); i++ {
		nextC := c.Clone()
		collectors = append(collectors, nextC)
	}

	ctxCtl, cancel := context.WithCancel(context.Background())
	ctxs := make([]*Context, 0, nodesLen)
	for i := 0; i < nodesLen; i++ {
		var ctx *Context
		if i != nodesLen-1 {
			ctx, err = newContext(ctxCtl, cancel, s.task, collectors[i], collectors[i+1])
		} else {
			ctx, err = newContext(ctxCtl, cancel, s.task, collectors[i], nil)
		}
		if err != nil {
			return err
		}
		ctxs = append(ctxs, ctx)
		if s.task.OutputConfig.Type == common.OutputTypeMySQL {
			if s.db != nil {
				ctx.setOutputDB(s.db)
			} else {
				db, err := common.NewDB(s.task.OutputConfig.MySQLConf)
				if err != nil {
					return err
				}
				ctx.setOutputDB(db)
			}
		}

		addCallback(ctx, s.task.Rule.Nodes[i])
	}

	headCtx, err := newContext(ctxCtl, cancel, s.task, c, collectors[0])
	if err != nil {
		return err
	}
	if s.task.OutputConfig.Type == common.OutputTypeMySQL {
		headCtx.setOutputDB(s.db)
	}
	headWrapper := func(ctx *Context) (err error) {
		defer func() {
			if e := recover(); e != nil {
				if v, ok := e.(error); ok {
					err = v
				} else {
					err = fmt.Errorf("%v", e)
				}
				log.Errorf("Head unexcepted exited, err: %+v, stack:\n%s", e, string(debug.Stack()))
			}
		}()
		return s.task.Rule.Head(ctx)
	}
	if err := headWrapper(headCtx); err != nil {
		log.Errorf("exec rule head func err:%#v", err)
		return errors.WithStack(err)
	}
	if err := addTaskCtrl(s.task.ID, cancel); err != nil {
		return errors.Wrapf(err, "addTaskCtrl failed")
	}

	go func() {
		for i := 0; i < nodesLen; i++ {
			collectors[i].Wait()
			log.Infof("task:%s %d step completed...", s.task.Name, i+1)
		}

		CancelTask(s.task.ID)
		s.retCh <- common.MTS{ID: s.task.ID, Status: common.TaskStatusCompleted}

		for _, ctx := range ctxs {
			ctx.closeCSVFileIfNeeded()
			if ctx.outputDB != nil {
				ctx.outputDB.Close()
			}
		}

		log.Infof("task:%s run completed...", s.task.Name)
	}()

	return nil
}

func cbDefer(ctx *Context, info string) {
	if e := recover(); e != nil {
		log.Error(info + fmt.Sprintf(", err: %+v, stack:\n%s", e, string(debug.Stack())))
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
				log.Warnf("request has been canceled in OnRequest, url:%s", newCtx.GetRequest().URL.String())
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
				log.Warnf("request has been canceled in OnError, url:%s", newCtx.GetRequest().URL.String())
				return
			default:
			}

			err := node.OnError(newCtx, newResponse(res, newCtx), e)
			if err != nil {
				log.Warnf("node.OnError return err:%+v, request url:%s", err, res.Request.URL.String())
			}
		})
	}

	if node.OnResponse != nil {
		ctx.c.OnResponse(func(res *colly.Response) {
			defer cbDefer(ctx, fmt.Sprintf("OnResponse unexcepted exited, url:%s", res.Request.URL.String()))

			newCtx := ctx.cloneWithReq(res.Request)
			select {
			case <-newCtx.ctlCtx.Done():
				log.Warnf("request has been canceled in OnResponse, url:%s", newCtx.GetRequest().URL.String())
				return
			default:
			}

			err := node.OnResponse(newCtx, newResponse(res, newCtx))
			if err != nil {
				log.Warnf("node.OnResponse return err:%+v, request url:%s", err, res.Request.URL.String())
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
					log.Warnf("request has been canceled in OnHTML, url:%s", newCtx.GetRequest().URL.String())
					return
				default:
				}

				err := f(newCtx, newHTMLElement(el, newCtx))
				if err != nil {
					log.Warnf("node.OnHTML:%s return err:%+v, request url:%s", selector, err, el.Request.URL.String())
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
					log.Warnf("request has been canceled in OnXML, url:%s", newCtx.GetRequest().URL.String())
					return
				default:
				}

				err := f(newCtx, newXMLElement(el, newCtx))
				if err != nil {
					log.Warnf("node.OnXML:%s return err:%+v, request url:%s", selector, err, el.Request.URL.String())
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
				log.Warnf("request has been canceled in OnScraped, url:%s", newCtx.GetRequest().URL.String())
				return
			default:
			}

			err := node.OnScraped(newCtx, newResponse(res, newCtx))
			if err != nil {
				log.Warnf("node.OnScraped return err:%+v, request url:%s", err, res.Request.URL.String())
			}
		})
	}

}
