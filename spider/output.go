package spider

import (
	"encoding/csv"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/go-xorm/builder"

	"github.com/nange/gospider/common"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// Output output a row data
func (ctx *Context) Output(row map[int]interface{}, namespace ...string) error {
	var outputFields []string
	var ns string

	switch len(namespace) {
	case 0:
		outputFields = ctx.task.OutputFields
		ns = ctx.task.Namespace
	case 1:
		if !ctx.task.OutputToMultipleNamespace {
			return ErrOutputToMultipleTableDisabled
		}
		multConf, ok := ctx.task.MultipleNamespaceConf[namespace[0]]
		if !ok {
			return ErrMultConfNamespaceNotFound
		}
		outputFields = multConf.OutputFields
		ns = namespace[0]
	default:
		return ErrTooManyOutputNamespace
	}

	if err := ctx.checkOutput(row, outputFields); err != nil {
		log.Errorf("checkOutput failed! err:%+v, fields:%#v, row:%+v", err, outputFields, row)
		return err
	}
	log.Debugf("output row:%+v", row)

	switch ctx.task.OutputConfig.Type {
	case common.OutputTypeMySQL:
		if err := ctx.outputToDB(row, outputFields, ns); err != nil {
			return err
		}
	case common.OutputTypeCSV:
		if err := ctx.outputToCSV(row, outputFields, ns); err != nil {
			return err
		}
	case common.OutputTypeStdout:
		if err := ctx.outputToStdout(row, outputFields, ns); err != nil {
			return err
		}
	default:
		return ErrOutputTypeNotSupported
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
		if !strings.HasPrefix(field, "`") {
			field = fmt.Sprintf("`%s`", field)
		}
		data[field] = row[i]
	}

	if !strings.HasPrefix(table, "`") {
		table = fmt.Sprintf("`%s`", table)
	}
	cond, vals, err := builder.Insert(builder.Eq(data)).Into(table).ToSQL()
	if err != nil {
		log.Errorf("build insert sql failed! err [%s], namespace [%s], row [%+v]", err.Error(), table, row)
		return errors.WithStack(err)
	}

	if _, err := ctx.outputDB.Exec(cond, vals...); err != nil {
		log.Errorf("exec insert sql failed! err:%s, cond:%s, vals:%+v", err.Error(), cond, vals)
		return errors.WithStack(err)
	}

	return nil
}

func (ctx *Context) outputToStdout(row map[int]interface{}, outputFields []string, ns string) error {
	fmt.Printf("output row:%+v\n", row)
	return nil
}

func (ctx *Context) outputToCSV(row map[int]interface{}, outputFields []string, ns string) error {
	w := ctx.outputCSVFiles[ns]
	cw := csv.NewWriter(w)

	record := make([]string, 0, len(outputFields))
	for i := 0; i < len(outputFields); i++ {
		record = append(record, fmt.Sprintf("%v", row[i]))
	}

	return errors.WithStack(cw.Write(record))
}

func createCSVFileIfNeeded(csvdir, csvfile string, outputFields []string) error {
	outputPath := path.Join(csvdir, csvfile)
	if _, err := os.Stat(outputPath); err == nil {
		return nil
	}

	if err := os.MkdirAll(csvdir, os.ModePerm); err != nil {
		return errors.Wrapf(err, "make csv output dir failed")
	}
	f, err := os.OpenFile(outputPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return errors.Wrapf(err, "create csv output file failed")
	}
	defer f.Close()

	w := csv.NewWriter(f)
	if err := w.Write(outputFields); err != nil {
		return errors.Wrapf(err, "write csv head failed")
	}
	w.Flush()
	if err := w.Error(); err != nil {
		return errors.Wrapf(err, "write csv head failed")
	}

	return nil
}
