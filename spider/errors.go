package spider

import "github.com/pkg/errors"

var (
	// ErrTaskRuleNotExist is the error type for task rule not exist
	ErrTaskRuleNotExist = errors.New("task rule not exist")
	// ErrTaskRuleIsNil is the error thrown when a nil rule registered
	ErrTaskRuleIsNil = errors.New("task rule is nil")
	// ErrTaskRuleNameIsEmpty is the error thrown when the ruleName is empty
	ErrTaskRuleNameIsEmpty = errors.New("task rule name is empty")
	// ErrTaskRuleNameDuplicated is the error thrown if the rule name is duplicated
	ErrTaskRuleNameDuplicated = errors.New("task rule name is Duplicated")
	// ErrTaskRuleHeadIsNil is the error thrown if the rule's head is nil
	ErrTaskRuleHeadIsNil = errors.New("task rule head is nil")
	// ErrTaskRuleNodesLenInvalid is the error thrown if the rule's nodes len is invalid
	ErrTaskRuleNodesLenInvalid = errors.New("task rule nodes len is invalid")
	// ErrTaskRuleNodesKeyInvalid is the error thrown if the rule's key len is invalid
	ErrTaskRuleNodesKeyInvalid = errors.New("task rule nodes key should start from 0 and monotonically increasing")
	// ErrTaskRunningTimeout is the error type for task running timeout
	ErrTaskRunningTimeout = errors.New("task running timeout")
)

var (
	// ErrOutputFieldsNotMatchOutputRow is the error type for output fields not match out put row
	ErrOutputFieldsNotMatchOutputRow = errors.New("output fields not match out put row")
	// ErrTooManyOutputNamespace is the error type for for too many output namespace
	ErrTooManyOutputNamespace = errors.New("too many output namespace")
	// ErrOutputToMultipleTableDisabled is the error thrown if "OutputToMultipleTable" is false
	ErrOutputToMultipleTableDisabled = errors.New("output to multiple tables disabled")
	// ErrOutputTypeNotSupported is the error type for unkonow output type
	ErrOutputTypeNotSupported = errors.New("output type not supported")
	// ErrMultConfNamespaceNotFound is the error type for mult conf namespace not found
	ErrMultConfNamespaceNotFound = errors.New("mult conf namespace not found")
)
