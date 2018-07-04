package service

import (
	"github.com/jinzhu/gorm"
	"github.com/nange/gospider/common"
	"github.com/nange/gospider/web/core"
	"github.com/nange/gospider/web/model"
	"github.com/sirupsen/logrus"
)

// 检查任务状态
// 1. 将TaskStatusRunning标记为TaskStatusUnexceptedExited
// 2. 将TaskStatusCompleted并且是定时任务的状态重新启动
func CheckTask() {
	qs := model.NewTaskQuerySet(core.GetDB())

	tasks := make([]model.Task, 0)
	if err := qs.All(&tasks); err != nil {
		if err == gorm.ErrRecordNotFound {
			logrus.Infof("no task found, exit service.CheckTask method")
			return
		}
		logrus.Errorf("query task list err: %+v", err)
		return
	}

	for _, task := range tasks {
		if task.Status == common.TaskStatusRunning {
			logrus.Infof("set task status to TaskStatusUnexceptedExited, taskID:%v", task.ID)
			err := model.NewTaskQuerySet(core.GetDB()).IDEq(task.ID).
				GetUpdater().SetStatus(common.TaskStatusUnexceptedExited).Update()
			if err != nil {
				logrus.Errorf("update task status err: %+v", err)
				continue
			}
		}

		if task.Status == common.TaskStatusCompleted && task.CronSpec != "" {
			if err := RestartTask(task); err != nil {
				logrus.Errorf("restart task err:%+v", err)
				continue
			}
		}
	}

}

// TODO:
func RestartTask(task model.Task) error {
	logrus.Infof("restarting task, taskID:%v", task.ID)

	return nil
}
