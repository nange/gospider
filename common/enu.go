package common

import (
	"errors"
	"fmt"
	"strings"
)

type TaskStatus uint8

const (
	TaskStatusUnknown TaskStatus = iota
	TaskStatusRunning
	TaskStatusPaused
	TaskStatusStopped
	TaskStatusUnexceptedExited
	TaskStatusCompleted
	TaskStatusRunningTimeout
)

var tsMap = map[TaskStatus]string{
	TaskStatusUnknown:          "未知状态",
	TaskStatusRunning:          "运行中",
	TaskStatusPaused:           "暂停",
	TaskStatusStopped:          "停止",
	TaskStatusUnexceptedExited: "异常退出",
	TaskStatusCompleted:        "完成",
	TaskStatusRunningTimeout:   "运行超时",
}

var InvalidTaskStatus = errors.New("invalid task status")

func (ts TaskStatus) String() string {
	s, ok := tsMap[ts]
	if !ok {
		panic(fmt.Sprintf("unexcepted TaskStatus %d", ts))
	}
	return s
}

func ParseTaskStatusFromString(s string) (TaskStatus, error) {
	switch s {
	case "未知状态":
		return TaskStatusUnknown, nil
	case "运行中":
		return TaskStatusRunning, nil
	case "暂停":
		return TaskStatusPaused, nil
	case "停止":
		return TaskStatusStopped, nil
	case "异常退出":
		return TaskStatusUnexceptedExited, nil
	case "完成":
		return TaskStatusCompleted, nil
	case "运行超时":
		return TaskStatusRunningTimeout, nil
	}
	return TaskStatusUnknown, InvalidTaskStatus
}

func (ts TaskStatus) MarshalJSON() ([]byte, error) {
	return []byte("\"" + ts.String() + "\""), nil
}

func (ts *TaskStatus) UnmarshalJSON(data []byte) (err error) {
	s := strings.Trim(strings.ToUpper(string(data)), "\"")
	*ts, err = ParseTaskStatusFromString(s)
	return
}
