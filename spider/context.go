package spider

import (
	"database/sql"

	qb "github.com/didi/gendry/builder"
	"github.com/gocolly/colly"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Context struct {
	task  *Task
	c     *colly.Collector
	nextC *colly.Collector

	// output
	outputDB *sql.DB
}

func newContext(task *Task, c *colly.Collector, nextC *colly.Collector) *Context {
	return &Context{
		task:  task,
		c:     c,
		nextC: nextC,
	}
}

func (ctx *Context) setOutputDB(db *sql.DB) {
	ctx.outputDB = db
}

func (ctx *Context) Visit(URL string) error {
	return ctx.c.Visit(URL)
}

func (ctx *Context) VisitForNext(URL string) error {
	return ctx.nextC.Visit(URL)
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

func (ctx *Context) PostMultipartForNext(URL string, requestData map[string][]byte) error {
	return ctx.nextC.PostMultipart(URL, requestData)
}

func (ctx *Context) Output(row map[int]interface{}) error {
	if err := ctx.checkOutput(row); err != nil {
		logrus.Errorf("checkOutput failed! err:%#v, fields:%#v, row:%#v", err, ctx.task.OutputFields, row)
		return err
	}
	logrus.Infof("output row:%#v", row)

	if ctx.task.OutputConfig.Type == OutputTypeMySQL {
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
		logrus.Errorf("build insert sql failed! err:%s, namespace:%s, row:%#v", err.Error(), ctx.task.Namespace, row)
		return errors.WithStack(err)
	}

	if _, err := ctx.outputDB.Exec(cond, vals...); err != nil {
		logrus.Errorf("exec insert sql failed! err:%s, cond:%s, vals:%#v", err.Error(), cond, vals)
		return errors.WithStack(err)
	}

	return nil
}

func (ctx *Context) outputToCVS(row map[int]interface{}) error {
	return nil
}
