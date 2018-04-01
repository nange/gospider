package spider

var rules = make(map[string]*TaskRule)

func Register(rule *TaskRule) {
	if err := checkRule(rule); err != nil {
		panic(err)
	}
	rules[rule.Name] = rule
}

type TaskRule struct {
	Name         string
	Description  string
	Namespace    string
	OutputFields []string
	Rule         *Rule
}

type Rule struct {
	Head  func(ctx *Context) error
	Nodes map[int]*Node
}

type Node struct {
	OnRequest  func(ctx *Context, req *Request)
	OnError    func(ctx *Context, res *Response, err error)
	OnResponse func(ctx *Context, res *Response)
	OnHTML     map[string]func(ctx *Context, el *HTMLElement)
	OnXML      map[string]func(ctx *Context, el *XMLElement)
	OnScraped  func(ctx *Context, res *Response)
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
