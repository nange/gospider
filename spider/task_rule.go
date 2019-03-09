package spider

import "github.com/pkg/errors"

var rules = make(map[string]*TaskRule)

// Register register a task rule
func Register(rule *TaskRule) {
	if err := checkRule(rule); err != nil {
		panic(err)
	}
	rules[rule.Name] = rule
}

// MultipleNamespaceConf is the mutiple namespace conf
type MultipleNamespaceConf struct {
	OutputFields      []string
	OutputConstraints map[string]*OutputConstraint
	OutputTableOpts   string
}

// TaskRule is the task rule define
type TaskRule struct {
	Name                      string
	Description               string
	OutputToMultipleNamespace bool
	MultipleNamespaceConf     map[string]*MultipleNamespaceConf
	Namespace                 string
	OutputFields              []string
	OutputConstraints         map[string]*OutputConstraint
	OutputTableOpts           string
	DisableCookies            bool
	AllowURLRevisit           bool
	IgnoreRobotsTxt           bool
	InsecureSkipVerify        bool
	ParseHTTPErrorResponse    bool
	Rule                      *Rule
}

// GetTaskRule get task rule by ruleName
func GetTaskRule(ruleName string) (*TaskRule, error) {
	if rule, ok := rules[ruleName]; ok {
		return rule, nil
	}
	return nil, errors.WithStack(ErrTaskRuleNotExist)
}

// GetTaskRuleKeys return all keys of task rule
func GetTaskRuleKeys() []string {
	keys := make([]string, 0, len(rules))
	for k := range rules {
		keys = append(keys, k)
	}

	return keys
}

// Rule the rule define
type Rule struct {
	Head  func(ctx *Context) error
	Nodes map[int]*Node
}

// Node the rule node of a task
type Node struct {
	OnRequest  func(ctx *Context, req *Request)
	OnError    func(ctx *Context, res *Response, err error) error
	OnResponse func(ctx *Context, res *Response) error
	OnHTML     map[string]func(ctx *Context, el *HTMLElement) error
	OnXML      map[string]func(ctx *Context, el *XMLElement) error
	OnScraped  func(ctx *Context, res *Response) error
}

func checkRule(rule *TaskRule) error {
	if rule == nil || rule.Rule == nil {
		return ErrTaskRuleIsNil
	}
	if rule.Name == "" {
		return ErrTaskRuleNameIsEmpty
	}
	if rule.Rule.Head == nil {
		return ErrTaskRuleHeadIsNil
	}
	if len(rule.Rule.Nodes) == 0 {
		return ErrTaskRuleNodesLenInvalid
	}
	for i := 0; i < len(rule.Rule.Nodes); i++ {
		if _, ok := rule.Rule.Nodes[i]; !ok {
			return ErrTaskRuleNodesKeyInvalid
		}
	}
	if _, ok := rules[rule.Name]; ok {
		return ErrTaskRuleNameDuplicated
	}

	return nil
}
