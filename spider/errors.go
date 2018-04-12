package spider

import "github.com/pkg/errors"

var (
	ErrTaskRuleNotExist        = errors.New("task rule not exist")
	ErrTaskRuleIsNil           = errors.New("task rule is nil")
	ErrTaskRuleNameIsEmpty     = errors.New("task rule name is empty")
	ErrTaskRuleNameDuplicated  = errors.New("task rule name is Duplicated")
	ErrTaskRuleHeadIsNil       = errors.New("task rule head is nil")
	ErrTaskRuleNodesLenInvalid = errors.New("task rule nodes len is invalid")
	ErrTaskRuleNodesKeyInvalid = errors.New("task rule nodes key should start from 0 and monotonically increasing")
)

var (
	ErrTaskRunningTimeout = errors.New("task running timeout")
)

var (
	ErrOutputFieldsNotMatchOutputRow = errors.New("output fields not match out put row")
)

var (
	ErrCronTaskDuplicated = errors.New("cron task is Duplicated")
)
