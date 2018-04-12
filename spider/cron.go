package spider

import (
	"sync"

	"github.com/nange/gospider/common"
	"github.com/pkg/errors"
	"github.com/robfig/cron"

	"github.com/sirupsen/logrus"
)

var cronTaskMap = &sync.Map{}

type CronTask struct {
	task  *Task
	cr    *cron.Cron
	retCh chan<- common.TaskStatus
}

func NewCronTask(task *Task, retCh chan<- common.TaskStatus) (*CronTask, error) {
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
	logrus.Infof("test run....")
	err := Run(ct.task, ct.retCh)
	if err != nil {
		logrus.Errorf("cron task run err:%+v", err)
	}
	ct.retCh <- common.TaskStatusRunning
}

func (ct *CronTask) Start() error {
	if err := ct.cr.AddJob(ct.task.CronSpec, ct); err != nil {
		return errors.Wrapf(err, "cron add job failed, task name:%s", ct.task.Name)
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
	if _, loaded := cronTaskMap.LoadOrStore(ct.task.Name, ct); loaded {
		return errors.Wrap(ErrCronTaskDuplicated, "add cron task failed, name:"+ct.task.Name)
	}
	return nil
}
