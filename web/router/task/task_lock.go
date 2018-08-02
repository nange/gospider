package task

import (
	"sync"
)

var taskLock *TaskLock

type TaskLock struct {
	taskLock   *sync.Map
	taskMaxLen uint64
}

func init() {
	taskLock = &TaskLock{
		taskLock:   &sync.Map{},
		taskMaxLen: uint64(1000),
	}
}

func (tl *TaskLock) IsLock(taskid uint64) bool {
	value, _ := tl.taskLock.LoadOrStore(taskid%tl.taskMaxLen, false)
	return value == true
}

func (tl *TaskLock) Lock(taskid uint64) {
	tl.taskLock.Store(taskid%tl.taskMaxLen, true)
}

func (tl *TaskLock) UnLock(taskid uint64) {
	tl.taskLock.Store(taskid%tl.taskMaxLen, false)
}
