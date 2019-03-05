package service

import (
	"github.com/jinzhu/gorm"
	"github.com/nange/gospider/common"
	"github.com/nange/gospider/web/core"
	"github.com/nange/gospider/web/model"
	log "github.com/sirupsen/logrus"
)

// check task status
// 1. set status to TaskStatusRunning if status is TaskStatusUnexceptedExited
// 2. restart task if status is TaskStatusCompleted, TaskStatusRunning
func CheckTask() {
	log.Infof("starting check task goroutine")
	qs := model.NewTaskQuerySet(core.GetGormDB())

	tasks := make([]model.Task, 0)
	if err := qs.All(&tasks); err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Infof("no task found, exit service.CheckTask method")
			return
		}
		log.Errorf("query task list err: %+v", err)
		return
	}

	for _, task := range tasks {
		if task.Status == common.TaskStatusRunning {
			log.Infof("set task status to TaskStatusUnexceptedExited, taskID:%v", task.ID)
			err := model.NewTaskQuerySet(core.GetGormDB()).IDEq(task.ID).
				GetUpdater().SetStatus(common.TaskStatusUnexceptedExited).Update()
			if err != nil {
				log.Errorf("update task status err: %+v", err)
				continue
			}
		}

		if (task.Status == common.TaskStatusCompleted || task.Status == common.TaskStatusRunning ||
			task.Status == common.TaskStatusUnexceptedExited) && task.CronSpec != "" {
			if err := CreateCronTask(task); err != nil {
				log.Errorf("restart task err:%+v", err)
				continue
			}
		}
	}

}

func CreateCronTask(task model.Task) error {
	log.Infof("restarting task, taskID:%v", task.ID)

	ct, err := NewCronTask(task.ID, task.CronSpec, GetMTSChan())
	if err != nil {
		log.Errorf("new cron task failed! err:%+v", err)
		return err
	}

	if err := ct.Start(); err != nil {
		log.Errorf("start cron task failed! err:%+v", err)
		return err
	}

	return nil
}
