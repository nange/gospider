package spider

import (
	"context"
	"fmt"
	"sync"

	"github.com/pkg/errors"
)

// Task is a task define
type Task struct {
	ID uint64
	TaskRule
	TaskConfig
}

// NewTask return a new task object
func NewTask(id uint64, rule TaskRule, config TaskConfig) *Task {
	return &Task{
		ID:         id,
		TaskRule:   rule,
		TaskConfig: config,
	}
}

var (
	ctlMu  = &sync.RWMutex{}
	ctlMap = make(map[uint64]context.CancelFunc)
)

func addTaskCtrl(taskID uint64, cancelFunc context.CancelFunc) error {
	ctlMu.Lock()
	defer ctlMu.Unlock()

	if _, ok := ctlMap[taskID]; ok {
		return errors.WithStack(fmt.Errorf("duplicate taskID:%d", taskID))
	}
	ctlMap[taskID] = cancelFunc

	return nil
}

// CancelTask cancel a task by taskID
func CancelTask(taskID uint64) bool {
	ctlMu.Lock()
	defer ctlMu.Unlock()

	cancel, ok := ctlMap[taskID]
	if !ok {
		return false
	}
	cancel()
	delete(ctlMap, taskID)

	return true
}
