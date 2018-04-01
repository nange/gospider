package spider

import (
	"database/sql"

	"github.com/gocolly/colly"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// TODO: Context添加KV功能，能够结束请求链功能
// TODO: 思考出错, 中断后续爬虫的方法
func Run(task *Task) error {
	var db *sql.DB
	var err error
	if task.OutputConfig.Type == OutputTypeMySQL {
		db, err = newDB(task.TaskConfig.OutputConfig.MySQLConf)
		if err != nil {
			logrus.Errorf("newDB failed! err:%#v", err)
			return err
		}
	}
	c := newCollector(task.TaskConfig)

	nodesLen := len(task.Rule.Nodes)
	cNodes := make([]*colly.Collector, 0, nodesLen)
	for i := 0; i < len(task.Rule.Nodes); i++ {
		cNodes = append(cNodes, c.Clone())
	}

	headCtx := newContext(task, c, cNodes[0])
	if err := task.Rule.Head(headCtx); err != nil {
		logrus.Errorf("exec rule head func err:%#v", err)
		return errors.WithStack(err)
	}

	var ctx *Context
	for i := 0; i < nodesLen; i++ {
		if i != nodesLen-1 {
			ctx = newContext(task, cNodes[i], cNodes[i+1])
		} else {
			ctx = newContext(task, cNodes[i], nil)
		}
		if task.OutputConfig.Type == OutputTypeMySQL {
			ctx.setOutputDB(db)
		}

		addCallback(ctx, task.Rule.Nodes[i])
	}
	for i := 0; i < nodesLen; i++ {
		cNodes[i].Wait()
	}
	logrus.Infof("task run completed...")
	return nil
}

func addCallback(ctx *Context, node *Node) {
	if node.OnRequest != nil {
		ctx.c.OnRequest(func(req *colly.Request) {
			node.OnRequest(ctx, newRequest(req, ctx))
		})
	}

	if node.OnError != nil {
		ctx.c.OnError(func(res *colly.Response, e error) {
			node.OnError(ctx, newResponse(res, ctx), e)
		})
	}

	if node.OnResponse != nil {
		ctx.c.OnResponse(func(res *colly.Response) {
			node.OnResponse(ctx, newResponse(res, ctx))
		})
	}

	if node.OnHTML != nil {
		for selector, fn := range node.OnHTML {
			ctx.c.OnHTML(selector, func(el *colly.HTMLElement) {
				fn(ctx, newHTMLElement(el, ctx))
			})
		}
	}

	if node.OnXML != nil {
		for selector, fn := range node.OnXML {
			ctx.c.OnXML(selector, func(el *colly.XMLElement) {
				fn(ctx, newXMLElement(el, ctx))
			})
		}
	}

	if node.OnScraped != nil {
		ctx.c.OnScraped(func(res *colly.Response) {
			node.OnScraped(ctx, newResponse(res, ctx))
		})
	}

}
