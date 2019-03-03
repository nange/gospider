package spider

import (
	"context"
	"os"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/nange/gospider/common"
	"github.com/stretchr/testify/suite"
)

type testOutputSuite struct {
	suite.Suite
	baseTask Task
}

func (s *testOutputSuite) SetupSuite() {
	s.baseTask = Task{
		ID:       1,
		TaskRule: TaskRule{Name: "test_task", Namespace: "test_table", OutputFields: []string{"field1", "field2"}},
		TaskConfig: TaskConfig{
			OutputConfig: OutputConfig{Type: common.OutputTypeMySQL},
		},
	}
}

func (s *testOutputSuite) TearDownSuite() {
	if s.DirExists("./csv_output") {
		s.NoError(os.RemoveAll("./csv_output"))
	}
}

func (s *testOutputSuite) TestDBOutputNormal() {
	db, mock, err := sqlmock.New()
	s.Require().NoError(err)
	defer db.Close()

	task := s.baseTask
	ctx, cancel := context.WithCancel(context.Background())
	gsCtx := newContext(ctx, cancel, &task, nil, nil)
	gsCtx.setOutputDB(db)

	mock.ExpectExec("(?i)insert into `test_table` (.+) values").
		WillReturnResult(sqlmock.NewResult(1, 1))

	row := map[int]interface{}{
		0: "field_value1",
		1: "field_value2",
	}
	err = gsCtx.Output(row)
	s.NoError(err)

	err = gsCtx.Output(row, "test_table")
	s.Equal(ErrOutputToMultipleTableDisabled, err)

	err = mock.ExpectationsWereMet()
	s.NoError(err)
}

func (s *testOutputSuite) TestDBOutputMult() {
	db, mock, err := sqlmock.New()
	s.Require().NoError(err)
	defer db.Close()

	task := s.baseTask
	task.TaskRule.OutputToMultipleNamespace = true
	task.TaskRule.MultipleNamespaceConf = map[string]*MultipleNamespaceConf{
		"test_mult_1": &MultipleNamespaceConf{OutputFields: []string{"field1", "field2"}},
		"test_mult_2": &MultipleNamespaceConf{OutputFields: []string{"field1", "field2"}},
	}
	ctx, cancel := context.WithCancel(context.Background())
	gsCtx := newContext(ctx, cancel, &task, nil, nil)
	gsCtx.setOutputDB(db)

	mock.ExpectExec("(?i)insert into `test_mult_1` (.+) values").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("(?i)insert into `test_mult_2` (.+) values").
		WillReturnResult(sqlmock.NewResult(2, 1))
	row := map[int]interface{}{
		0: "field_value1",
		1: "field_value2",
	}
	err = gsCtx.Output(row, "test_mult_1")
	s.NoError(err)
	err = gsCtx.Output(row, "test_mult_2")
	s.NoError(err)

	err = mock.ExpectationsWereMet()
	s.NoError(err)

	err = gsCtx.Output(row, "test_mult_not_exist")
	s.Equal(err, ErrMultConfNamespaceNotFound)

	err = gsCtx.Output(row, "args_too_much", "args_too_much2")
	s.Equal(err, ErrTooManyOutputNamespace)
}

func (s *testOutputSuite) TestCSVOutputNormal() {
	task := s.baseTask
	task.TaskConfig.OutputConfig.Type = common.OutputTypeCSV
	task.TaskConfig.OutputConfig.CSVConf.CSVFilePath = "./csv_output"

	ctx, cancel := context.WithCancel(context.Background())
	gsCtx := newContext(ctx, cancel, &task, nil, nil)

	row := map[int]interface{}{
		0: "field_value1",
		1: "field_value2",
	}
	err := gsCtx.Output(row)
	s.NoError(err)

	gsCtx.ctlCancel()
	fileExist := false
	for i := 0; i < 3; i++ {
		if !fileExist {
			time.Sleep(500 * time.Millisecond)
		}
		fileExist = s.FileExists("./csv_output/test_table.csv")
		if fileExist {
			break
		}
	}
	s.True(fileExist)

	// mult csv output
	task.TaskRule.OutputToMultipleNamespace = true
	task.TaskRule.MultipleNamespaceConf = map[string]*MultipleNamespaceConf{
		"test_mult_1": &MultipleNamespaceConf{OutputFields: []string{"field1", "field2"}},
		"test_mult_2": &MultipleNamespaceConf{OutputFields: []string{"field1", "field2"}},
	}
	ctx, cancel = context.WithCancel(context.Background())
	gsCtx2 := newContext(ctx, cancel, &task, nil, nil)

	err = gsCtx2.Output(row, "test_mult_1")
	s.NoError(err)
	err = gsCtx2.Output(row, "test_mult_2")
	s.NoError(err)

}

func TestOutputSuite(t *testing.T) {
	suite.Run(t, new(testOutputSuite))
}
