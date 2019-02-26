package spider

import (
	"encoding/csv"
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"

	qb "github.com/didi/gendry/builder"
	"github.com/nange/gospider/common"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type outputRecord struct {
	Namespace    string
	OutputType   string
	OutputFields []string
	OutputItems  []interface{}
}

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
		outputFields = ctx.task.MultipleNamespaceConf[namespace[0]].OutputFields
		ns = namespace[0]
	default:
		return ErrTooManyOutputTables
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
		od := &outputRecord{
			Namespace:    ns,
			OutputType:   common.OutputTypeCSV,
			OutputFields: outputFields,
		}
		outputItems := make([]interface{}, 0, len(outputFields))
		for i := 0; i < len(outputFields); i++ {
			outputItems = append(outputItems, row[i])
		}
		od.OutputItems = outputItems
		select {
		case ctx.outputChan <- od:
			log.Debugf("write csv row to memory success")
		default:
			log.Errorf("write csv row [%+v] to chan failed", row)
			return errors.New("write csv row to chan failed")
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
		data[field] = row[i]
	}

	cond, vals, err := qb.BuildInsert(table, []map[string]interface{}{data})
	if err != nil {
		log.Errorf("build insert sql failed! err:%s, namespace:%s, row:%+v", err.Error(), table, row)
		return errors.WithStack(err)
	}

	quotedCond, err := quoteQuery(cond)
	if err != nil {
		log.Error(err)
		return errors.WithStack(err)
	}

	if _, err := ctx.outputDB.Exec(quotedCond, vals...); err != nil {
		log.Errorf("exec insert sql failed! err:%s, cond:%s, vals:%+v", err.Error(), quotedCond, vals)
		return errors.WithStack(err)
	}

	return nil
}

func (ctx *Context) asyncWriteCSVFile() {
	var csvfile *os.File
	var csvWriter *csv.Writer

	defer func() {
		if csvfile != nil {
			csvfile.Close()
		}
	}()

	for {
		select {
		case <-ctx.ctlCtx.Done():
			log.Debugf("task context done, taskID [%v]", ctx.task.ID)
			if csvWriter != nil {
				csvWriter.Flush()
			}
			return
		case record := <-ctx.outputChan:
			if csvfile == nil {
				csvConf := ctx.task.OutputConfig.CSVConf
				csvname := fmt.Sprintf("%s.csv", record.Namespace)
				err := createCSVFileIfNeeded(csvConf.CSVFilePath, csvname, record.OutputFields)
				if err != nil {
					log.Errorf("createCSVFileIfNeeded err [%+v], csvname [%v], we should cancel the task", err, csvname)
					CancelTask(ctx.task.ID)
					return
				}
				outputPath := path.Join(csvConf.CSVFilePath, csvname)
				csvfile, err = os.OpenFile(outputPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
				if err != nil {
					log.Errorf("open csv file [%s] err [%+v], we should cancel the task",
						outputPath, errors.WithStack(err))
					CancelTask(ctx.task.ID)
					return
				}
				csvWriter = csv.NewWriter(csvfile)
			}

			strItems := make([]string, 0, len(record.OutputItems))
			for _, item := range record.OutputItems {
				strItems = append(strItems, fmt.Sprintf("%v", item))
			}
			if err := csvWriter.Write(strItems); err != nil {
				log.Errorf("write csv record err [%+v]", errors.WithStack(err))
				break
			}
		}
	}
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
