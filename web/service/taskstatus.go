package service

import (
	"github.com/nange/gospider/common"
	"github.com/nange/gospider/web/core"
	"github.com/nange/gospider/web/model"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var mtsCh = make(chan common.MTS, 1)

func GetMTSChan() chan common.MTS {
	return mtsCh
}

func ManageTaskStatus() {
	logrus.Infof("starting manage task status goroutine")
	for {
		select {
		case mts := <-mtsCh:
			task := &model.Task{}
			err := model.NewTaskQuerySet(core.GetDB()).IDEq(mts.ID).One(task)
			if err != nil {
				logrus.Errorf("query model task err: %+v", err)
				break
			}

			task.Status = mts.Status
			if mts.Status == common.TaskStatusCompleted {
				task.Counts += 1
			}

			if err := task.Update(core.GetDB(), model.TaskDBSchema.Status, model.TaskDBSchema.Counts); err != nil {
				logrus.Errorf("update task status err:%+v", errors.WithStack(err))
				break
			}

		}
	}
}
