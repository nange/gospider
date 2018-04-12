package common

import "fmt"

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

func (ts TaskStatus) String() string {
	s, ok := tsMap[ts]
	if !ok {
		panic(fmt.Sprintf("unexcepted TaskStatus %d", ts))
	}
	return s
}

func (ts TaskStatus) MarshalJSON() ([]byte, error) {
	return []byte("\"" + ts.String() + "\""), nil
}
