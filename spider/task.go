package spider

type Task struct {
	TaskRule
	TaskConfig
}

func NewTask(rule TaskRule, config TaskConfig) *Task {
	return &Task{
		TaskRule: rule,
		TaskConfig: config,
	}
}
