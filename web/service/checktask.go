package service

import (
	"github.com/jinzhu/gorm"
	"github.com/nange/gospider/common"
	"github.com/nange/gospider/web/core"
	"github.com/nange/gospider/web/model"
	"github.com/sirupsen/logrus"
)

// check task status
// 1. set status to TaskStatusRunning if status is TaskStatusUnexceptedExited
// 2. restart task if status is TaskStatusCompleted, TaskStatusRunning
func CheckTask() {
	logrus.Infof("starting check task goroutine")
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

		if (task.Status == common.TaskStatusCompleted || task.Status == common.TaskStatusRunning ||
			task.Status == common.TaskStatusUnexceptedExited) && task.CronSpec != "" {
			if err := CreateCronTask(task); err != nil {
				logrus.Errorf("restart task err:%+v", err)
				continue
			}
		}
	}

}

func CreateCronTask(task model.Task) error {
	logrus.Infof("restarting task, taskID:%v", task.ID)

	ct, err := NewCronTask(task.ID, task.CronSpec, GetMTSChan())
	if err != nil {
		logrus.Errorf("new cron task failed! err:%+v", err)
		return err
	}

	if err := ct.Start(); err != nil {
		logrus.Errorf("start cron task failed! err:%+v", err)
		return err
	}

	return nil
}
