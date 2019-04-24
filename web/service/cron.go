package service

import (
	"fmt"
	"sync"

	"github.com/nange/gospider/common"
	"github.com/nange/gospider/spider"
	"github.com/nange/gospider/web/core"
	"github.com/nange/gospider/web/model"
	"github.com/pkg/errors"
	"github.com/robfig/cron"
	log "github.com/sirupsen/logrus"
)

var (
	ErrCronTaskDuplicated = errors.New("cron task is Duplicated")
)

var cronTaskMap = &sync.Map{}

type CronTask struct {
	taskID   uint64
	cronSpec string
	cr       *cron.Cron
	retCh    chan<- common.MTS
}

func NewCronTask(taskID uint64, cronSpec string, retCh chan<- common.MTS) (*CronTask, error) {
	ct := &CronTask{
		taskID:   taskID,
		cronSpec: cronSpec,
		cr:       cron.New(),
		retCh:    retCh,
	}

	if err := AddCronTask(ct); err != nil {
		return nil, errors.Wrap(err, "add cron task failed")
	}

	return ct, nil
}

func (ct *CronTask) Run() {
	task := &model.Task{}
	err := model.NewTaskQuerySet(core.GetGormDB()).IDEq(ct.taskID).One(task)
	if err != nil {
		log.Errorf("run cron task failed, query task err:%+v", errors.WithStack(err))
		return
	}
	if task.Status != common.TaskStatusCompleted && task.Status != common.TaskStatusUnexceptedExited {
		log.Warnf("run cron task failed, status:%+v", errors.New(task.Status.String()))
		return
	}

	spiderTask, err := GetSpiderTaskByModel(task)
	if err != nil {
		log.Errorf("run cron task failed, err:%+v", errors.WithStack(err))
		return
	}
	s := spider.New(spiderTask, ct.retCh)
	if err := s.Run(); err != nil {
		log.Errorf("run cron task failed, err:%+v", err)
		ct.retCh <- common.MTS{ID: task.ID, Status: common.TaskStatusUnexceptedExited}
		return
	}

	ct.retCh <- common.MTS{ID: task.ID, Status: common.TaskStatusRunning}
}

func (ct *CronTask) Start() error {
	if err := ct.cr.AddJob(ct.cronSpec, ct); err != nil {
		return errors.Wrapf(err, "cron add job failed, taskID:%d", ct.taskID)
	}
	ct.cr.Start()
	return nil
}

func (ct *CronTask) Stop() error {
	if ct.cr == nil {
		return errors.New("CronTask do not started")
	}
	cronTaskMap.Delete(ct.taskID)

	ct.cr.Stop()
	return nil
}

func GetCronTask(taskID uint64) *CronTask {
	ct, ok := cronTaskMap.Load(taskID)
	if !ok {
		return nil
	}
	return ct.(*CronTask)
}

func AddCronTask(ct *CronTask) error {
	if _, loaded := cronTaskMap.LoadOrStore(ct.taskID, ct); loaded {
		return errors.Wrap(ErrCronTaskDuplicated, fmt.Sprintf("add cron task failed, taskID:%d", ct.taskID))
	}
	return nil
}
