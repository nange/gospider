package spider

import (
	"context"
	"database/sql"
	"io"
	"net/http"

	qb "github.com/didi/gendry/builder"
	"github.com/gocolly/colly"
	"github.com/nange/gospider/common"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Context struct {
	task  *Task
	c     *colly.Collector
	nextC *colly.Collector

	ctlCtx context.Context

	// output
	outputDB *sql.DB
}

func newContext(ctx context.Context, task *Task, c *colly.Collector, nextC *colly.Collector) *Context {
	return &Context{
		task:   task,
		c:      c,
		nextC:  nextC,
		ctlCtx: ctx,
	}
}

func (ctx *Context) cloneWithReq(req *colly.Request) *Context {
	newctx := context.WithValue(ctx.ctlCtx, "req", req)

	return &Context{
		task:     ctx.task,
		c:        ctx.c,
		nextC:    ctx.nextC,
		ctlCtx:   newctx,
		outputDB: ctx.outputDB,
	}
}

func (ctx *Context) setOutputDB(db *sql.DB) {
	ctx.outputDB = db
}

func (ctx *Context) GetRequest() *Request {
	collyReq := ctx.ctlCtx.Value("req").(*colly.Request)
	return newRequest(collyReq, ctx)
}

func (ctx *Context) Retry() error {
	return ctx.ctlCtx.Value("req").(*colly.Request).Retry()
}

func (ctx *Context) PutReqContextValue(key string, value interface{}) {
	ctx.ctlCtx.Value("req").(*colly.Request).Ctx.Put(key, value)
}

func (ctx *Context) GetReqContextValue(key string) string {
	return ctx.ctlCtx.Value("req").(*colly.Request).Ctx.Get(key)
}

func (ctx *Context) GetAnyReqContextValue(key string) interface{} {
	return ctx.ctlCtx.Value("req").(*colly.Request).Ctx.GetAny(key)
}

func (ctx *Context) Visit(URL string) error {
	if req, ok := ctx.ctlCtx.Value("req").(*colly.Request); ok {
		return ctx.c.Visit(req.AbsoluteURL(URL))
	}
	return ctx.c.Visit(URL)
}

func (ctx *Context) VisitForNext(URL string) error {
	if req, ok := ctx.ctlCtx.Value("req").(*colly.Request); ok {
		return ctx.nextC.Visit(req.AbsoluteURL(URL))
	}
	return ctx.nextC.Visit(URL)
}

func (ctx *Context) reqContextClone() *colly.Context {
	newCtx := colly.NewContext()
	req := ctx.ctlCtx.Value("req").(*colly.Request)
	req.Ctx.ForEach(func(k string, v interface{}) interface{} {
		newCtx.Put(k, v)
		return nil
	})

	return newCtx
}

func (ctx *Context) VisitForNextWithContext(URL string) error {
	req := ctx.ctlCtx.Value("req").(*colly.Request)
	return ctx.nextC.Request("GET", req.AbsoluteURL(URL), nil, ctx.reqContextClone(), nil)
}

func (ctx *Context) Post(URL string, requestData map[string]string) error {
	return ctx.c.Post(URL, requestData)
}

func (ctx *Context) PostForNext(URL string, requestData map[string]string) error {
	return ctx.nextC.Post(URL, requestData)
}

func (ctx *Context) PostRawForNext(URL string, requestData []byte) error {
	return ctx.nextC.PostRaw(URL, requestData)
}

func (ctx *Context) RequestForNext(method, URL string, requestData io.Reader, hdr http.Header) error {
	return ctx.nextC.Request(method, URL, requestData, nil, hdr)
}

func (ctx *Context) PostMultipartForNext(URL string, requestData map[string][]byte) error {
	return ctx.nextC.PostMultipart(URL, requestData)
}

func (ctx *Context) Output(row map[int]interface{}) error {
	if err := ctx.checkOutput(row); err != nil {
		logrus.Errorf("checkOutput failed! err:%+v, fields:%#v, row:%+v", err, ctx.task.OutputFields, row)
		return err
	}
	logrus.Infof("output row:%+v", row)

	if ctx.task.OutputConfig.Type == common.OutputTypeMySQL {
		if err := ctx.outputToDB(row); err != nil {
			return err
		}
	}

	return nil
}

func (ctx *Context) checkOutput(row map[int]interface{}) error {
	if len(ctx.task.OutputFields) != len(row) {
		return ErrOutputFieldsNotMatchOutputRow
	}

	for i := 0; i < len(ctx.task.OutputFields); i++ {
		if _, ok := row[i]; !ok {
			return ErrOutputFieldsNotMatchOutputRow
		}
	}

	return nil
}

func (ctx *Context) outputToDB(row map[int]interface{}) error {
	data := make(map[string]interface{})
	for i, field := range ctx.task.OutputFields {
		data[field] = row[i]
	}

	cond, vals, err := qb.BuildInsert(ctx.task.Namespace, []map[string]interface{}{data})
	if err != nil {
		logrus.Errorf("build insert sql failed! err:%s, namespace:%s, row:%+v", err.Error(), ctx.task.Namespace, row)
		return errors.WithStack(err)
	}

	if _, err := ctx.outputDB.Exec(cond, vals...); err != nil {
		logrus.Errorf("exec insert sql failed! err:%s, cond:%s, vals:%+v", err.Error(), cond, vals)
		return errors.WithStack(err)
	}

	return nil
}

func (ctx *Context) outputToCVS(row map[int]interface{}) error {
	return nil
}
