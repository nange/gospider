package spider

import (
	"bytes"
	"context"
	"database/sql"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/gocolly/colly"
	"github.com/nange/gospider/common"
)

type Context struct {
	task  *Task
	c     *colly.Collector
	nextC *colly.Collector

	ctlCtx    context.Context
	ctlCancel context.CancelFunc

	collyContext *colly.Context
	// output
	outputDB   *sql.DB
	outputChan chan *outputRecord
}

func newContext(ctx context.Context, cancel context.CancelFunc, task *Task, c *colly.Collector, nextC *colly.Collector) *Context {
	gsCtx := &Context{
		task:       task,
		c:          c,
		nextC:      nextC,
		ctlCtx:     ctx,
		ctlCancel:  cancel,
		outputChan: make(chan *outputRecord, 64),
	}

	if task.OutputConfig.Type == common.OutputTypeCSV {
		go gsCtx.asyncWriteCSVFile()
	}

	return gsCtx
}

func (ctx *Context) cloneWithReq(req *colly.Request) *Context {
	newctx := context.WithValue(ctx.ctlCtx, "req", req)

	return &Context{
		task:       ctx.task,
		c:          ctx.c,
		nextC:      ctx.nextC,
		ctlCtx:     newctx,
		ctlCancel:  ctx.ctlCancel,
		outputDB:   ctx.outputDB,
		outputChan: ctx.outputChan,
	}
}

func (ctx *Context) setOutputDB(db *sql.DB) {
	ctx.outputDB = db
}

func (ctx *Context) GetRequest() *Request {
	if req, ok := ctx.ctlCtx.Value("req").(*colly.Request); ok {
		return newRequest(req, ctx)
	}
	return nil
}

func (ctx *Context) Retry() error {
	if req, ok := ctx.ctlCtx.Value("req").(*colly.Request); ok {
		return req.Retry()
	}

	return nil
}

func (ctx *Context) PutReqContextValue(key string, value interface{}) {
	if ctx.collyContext == nil {
		if req, ok := ctx.ctlCtx.Value("req").(*colly.Request); ok {
			ctx.collyContext = req.Ctx
		} else {
			ctx.collyContext = colly.NewContext()
		}
	}
	ctx.collyContext.Put(key, value)
}

func (ctx *Context) GetReqContextValue(key string) string {
	if ctx.collyContext == nil {
		if req, ok := ctx.ctlCtx.Value("req").(*colly.Request); ok {
			ctx.collyContext = req.Ctx
		} else {
			return ""
		}
	}
	return ctx.collyContext.Get(key)
}

func (ctx *Context) GetAnyReqContextValue(key string) interface{} {
	if ctx.collyContext == nil {
		if req, ok := ctx.ctlCtx.Value("req").(*colly.Request); ok {
			ctx.collyContext = req.Ctx
		} else {
			return nil
		}
	}
	return ctx.collyContext.GetAny(key)
}

func (ctx *Context) Visit(URL string) error {
	return ctx.c.Visit(ctx.AbsoluteURL(URL))
}

func (ctx *Context) VisitForNext(URL string) error {
	return ctx.nextC.Visit(ctx.AbsoluteURL(URL))
}

func (ctx *Context) reqContextClone() *colly.Context {
	newCtx := colly.NewContext()
	if ctx.collyContext == nil {
		return newCtx
	}

	ctx.collyContext.ForEach(func(k string, v interface{}) interface{} {
		newCtx.Put(k, v)
		return nil
	})

	return newCtx
}

func (ctx *Context) VisitForNextWithContext(URL string) error {
	return ctx.nextC.Request("GET", ctx.AbsoluteURL(URL), nil, ctx.reqContextClone(), nil)
}

func (ctx *Context) Post(URL string, requestData map[string]string) error {
	return ctx.c.Post(ctx.AbsoluteURL(URL), requestData)
}

func (ctx *Context) PostForNext(URL string, requestData map[string]string) error {
	return ctx.nextC.Post(ctx.AbsoluteURL(URL), requestData)
}

func (ctx *Context) PostForNextWithContext(URL string, requestData map[string]string) error {
	return ctx.nextC.Request("POST", ctx.AbsoluteURL(URL), createFormReader(requestData), ctx.reqContextClone(), nil)
}

func (ctx *Context) PostRawForNext(URL string, requestData []byte) error {
	return ctx.nextC.PostRaw(ctx.AbsoluteURL(URL), requestData)
}

func (ctx *Context) PostRawForNextWithContext(URL string, requestData []byte) error {
	return ctx.nextC.Request("POST", ctx.AbsoluteURL(URL), bytes.NewReader(requestData), ctx.reqContextClone(), nil)
}

func (ctx *Context) Request(method, URL string, requestData io.Reader, hdr http.Header) error {
	return ctx.c.Request(method, URL, requestData, nil, hdr)
}

func (ctx *Context) RequestWithContext(method, URL string, requestData io.Reader, hdr http.Header) error {
	return ctx.c.Request(method, URL, requestData, ctx.reqContextClone(), hdr)
}

func (ctx *Context) RequestForNext(method, URL string, requestData io.Reader, hdr http.Header) error {
	return ctx.nextC.Request(method, URL, requestData, nil, hdr)
}

func (ctx *Context) RequestForNextWithContext(method, URL string, requestData io.Reader, hdr http.Header) error {
	return ctx.nextC.Request(method, URL, requestData, ctx.reqContextClone(), hdr)
}

func (ctx *Context) PostMultipartForNext(URL string, requestData map[string][]byte) error {
	return ctx.nextC.PostMultipart(URL, requestData)
}

func (ctx *Context) SetResponseCharacterEncoding(encoding string) {
	if req, ok := ctx.ctlCtx.Value("req").(*colly.Request); ok {
		req.ResponseCharacterEncoding = encoding
	}
}

func (ctx *Context) AbsoluteURL(u string) string {
	if req, ok := ctx.ctlCtx.Value("req").(*colly.Request); ok {
		return req.AbsoluteURL(u)
	}
	return u
}

func (ctx *Context) Abort() {
	if req, ok := ctx.ctlCtx.Value("req").(*colly.Request); ok {
		req.Abort()
	}
}

func createFormReader(data map[string]string) io.Reader {
	form := url.Values{}
	for k, v := range data {
		form.Add(k, v)
	}
	return strings.NewReader(form.Encode())
}
