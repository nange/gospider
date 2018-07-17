package spider

import (
	"context"
	"fmt"
	"sync"

	"github.com/pkg/errors"
)

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

func CancelTask(taskID uint64) error {
	ctlMu.Lock()
	defer ctlMu.Unlock()

	cancel, ok := ctlMap[taskID]
	if !ok {
		return errors.WithStack(fmt.Errorf("taskID:%d not found", taskID))
	}
	cancel()
	delete(ctlMap, taskID)

	return nil
}
