package task

import (
	"sync"
)

var taskLock *TaskLock

type TaskLock struct {
	taskLock map[uint64]bool
	sync.Mutex
}

func init() {
	taskLock = &TaskLock{
		taskLock: make(map[uint64]bool),
	}
}

func (tl *TaskLock) IsRunning(taskid uint64) bool {
	tl.Lock()
	defer tl.Unlock()
	if tl.taskLock[taskid] {
		return true
	}
	tl.taskLock[taskid] = true
	return false
}

func (tl *TaskLock) Complete(taskid uint64) {
	tl.Lock()
	defer tl.Unlock()
	delete(tl.taskLock, taskid)
}
