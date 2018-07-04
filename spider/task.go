package spider

type Task struct {
	ID uint64
	TaskRule
	TaskConfig
}

func NewTask(id uint64, rule TaskRule, config TaskConfig) *Task {
	return &Task{
		ID:         id,
		TaskRule:   rule,
		TaskConfig: config,
	}
}
