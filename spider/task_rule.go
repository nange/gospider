package spider

import "github.com/pkg/errors"

var rules = make(map[string]*TaskRule)

func Register(rule *TaskRule) {
	if err := checkRule(rule); err != nil {
		panic(err)
	}
	rules[rule.Name] = rule
}

type MultipleNamespacesConf struct {
	OutputFields      []string
	OutputConstraints map[string]*OutputConstraint
	OutputTableOpts   string
}

type TaskRule struct {
	Name                       string
	Description                string
	OutputToMultipleNamespaces bool
	MultipleNamespacesConf     map[string]*MultipleNamespacesConf
	Namespace                  string
	OutputFields               []string
	OutputConstraints          map[string]*OutputConstraint
	OutputTableOpts            string
	DisableCookies             bool
	AllowURLRevisit            bool
	IgnoreRobotsTxt            bool
	InsecureSkipVerify         bool
	ParseHTTPErrorResponse     bool
	Rule                       *Rule
}

func GetTaskRule(ruleName string) (*TaskRule, error) {
	if rule, ok := rules[ruleName]; ok {
		return rule, nil
	}
	return nil, errors.WithStack(ErrTaskRuleNotExist)
}

func GetTaskRuleKeys() []string {
	keys := make([]string, 0, len(rules))
	for k := range rules {
		keys = append(keys, k)
	}

	return keys
}

type Rule struct {
	Head  func(ctx *Context) error
	Nodes map[int]*Node
}

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
