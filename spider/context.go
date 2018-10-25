package spider

import (
	"bytes"
	"context"
	"database/sql"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"

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

	ctlCtx    context.Context
	ctlCancel context.CancelFunc

	collyContext *colly.Context
	// output
	outputDB *sql.DB
}

func newContext(ctx context.Context, cancel context.CancelFunc, task *Task, c *colly.Collector, nextC *colly.Collector) *Context {
	return &Context{
		task:      task,
		c:         c,
		nextC:     nextC,
		ctlCtx:    ctx,
		ctlCancel: cancel,
	}
}

func (ctx *Context) cloneWithReq(req *colly.Request) *Context {
	newctx := context.WithValue(ctx.ctlCtx, "req", req)

	return &Context{
		task:      ctx.task,
		c:         ctx.c,
		nextC:     ctx.nextC,
		ctlCtx:    newctx,
		ctlCancel: ctx.ctlCancel,
		outputDB:  ctx.outputDB,
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

func (ctx *Context) Output(row map[int]interface{}, namespace ...string) error {
	var outputFields []string
	var ns string

	switch len(namespace) {
	case 0:
		outputFields = ctx.task.OutputFields
		ns = ctx.task.Namespace
	case 1:
		if !ctx.task.OutputToMultipleNamespaces {
			return ErrOutputToMultipleTableDisabled
		}
		outputFields = ctx.task.MultipleNamespacesConf[namespace[0]].OutputFields
		ns = namespace[0]
	default:
		return ErrTooManyOutputTables
	}

	if err := ctx.checkOutput(row, outputFields); err != nil {
		logrus.Errorf("checkOutput failed! err:%+v, fields:%#v, row:%+v", err, outputFields, row)
		return err
	}
	logrus.Infof("output row:%+v", row)

	if ctx.task.OutputConfig.Type == common.OutputTypeMySQL {
		if err := ctx.outputToDB(row, outputFields, ns); err != nil {
			return err
		}
	}

	return nil
}

func (ctx *Context) checkOutput(row map[int]interface{}, outputFields []string) error {
	if len(outputFields) != len(row) {
		return ErrOutputFieldsNotMatchOutputRow
	}

	for i := 0; i < len(outputFields); i++ {
		if _, ok := row[i]; !ok {
			return ErrOutputFieldsNotMatchOutputRow
		}
	}

	return nil
}

func (ctx *Context) outputToDB(row map[int]interface{}, outputFields []string, table string) error {
	data := make(map[string]interface{})
	for i, field := range outputFields {
		data[field] = row[i]
	}

	cond, vals, err := qb.BuildInsert(table, []map[string]interface{}{data})
	if err != nil {
		logrus.Errorf("build insert sql failed! err:%s, namespace:%s, row:%+v", err.Error(), table, row)
		return errors.WithStack(err)
	}

	quotedCond, err := quoteQuery(cond)
	if err != nil {
		logrus.Error(err)
		return errors.WithStack(err)
	}

	if _, err := ctx.outputDB.Exec(quotedCond, vals...); err != nil {
		logrus.Errorf("exec insert sql failed! err:%s, cond:%s, vals:%+v", err.Error(), quotedCond, vals)
		return errors.WithStack(err)
	}

	return nil
}

func (ctx *Context) outputToCVS(row map[int]interface{}) error {
	return nil
}

func quoteQuery(sql string) (s string, err error) {
	reg := regexp.MustCompile(`(?sU)(INSERT INTO .+ \(\s*)(.+)(\s*\) VALUES.+\))`)
	matches := reg.FindStringSubmatch(sql)
	if len(matches) != 4 {
		err = errors.New("quote sql regexp not match")
		return
	}
	fields := strings.Replace(matches[2], ",", "`,`", -1)
	s = matches[1] + "`" + fields + "`" + matches[3]
	return
}

func createFormReader(data map[string]string) io.Reader {
	form := url.Values{}
	for k, v := range data {
		form.Add(k, v)
	}
	return strings.NewReader(form.Encode())
}
