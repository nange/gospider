package service

import (
	"sync"

	"github.com/nange/gospider/common"
	"github.com/nange/gospider/spider"
	"github.com/nange/gospider/web/model"
	"github.com/pkg/errors"
	"github.com/robfig/cron"

	"github.com/sirupsen/logrus"
)

var (
	ErrCronTaskDuplicated = errors.New("cron task is Duplicated")
)

var cronTaskMap = &sync.Map{}

type CronTask struct {
	task  *model.Task
	cr    *cron.Cron
	retCh chan<- common.MTS
}

func NewCronTask(task *model.Task, retCh chan<- common.MTS) (*CronTask, error) {
	ct := &CronTask{
		task:  task,
		cr:    cron.New(),
		retCh: retCh,
	}

	if err := AddCronTask(ct); err != nil {
		return nil, errors.Wrap(err, "add cron task failed")
	}

	return ct, nil
}

func (ct *CronTask) Run() {
	spiderTask, err := GetSpiderTaskByModel(ct.task)
	if err != nil {
		logrus.Errorf("cron task run err:%+v", errors.WithStack(err))
		return
	}
	if err := spider.Run(spiderTask, ct.retCh); err != nil {
		logrus.Errorf("cron task run err:%+v", errors.WithStack(err))
		ct.retCh <- common.MTS{ID: ct.task.ID, Status: common.TaskStatusUnexceptedExited}
		return
	}

	ct.retCh <- common.MTS{ID: ct.task.ID, Status: common.TaskStatusRunning}
}

func (ct *CronTask) Start() error {
	if err := ct.cr.AddJob(ct.task.CronSpec, ct); err != nil {
		return errors.Wrapf(err, "cron add job failed, task name:%s", ct.task.TaskName)
	}
	ct.cr.Start()
	return nil
}

func (ct *CronTask) Stop() error {
	if ct.cr == nil {
		return errors.New("CronTask do not started")
	}

	ct.cr.Stop()
	return nil
}

func GetCronTask(name string) (*CronTask, bool) {
	ct, ok := cronTaskMap.Load(name)
	if !ok {
		return nil, false
	}
	return ct.(*CronTask), true
}

func AddCronTask(ct *CronTask) error {
	if _, loaded := cronTaskMap.LoadOrStore(ct.task.TaskName, ct); loaded {
		return errors.Wrap(ErrCronTaskDuplicated, "add cron task failed, name:"+ct.task.TaskName)
	}
	return nil
}
