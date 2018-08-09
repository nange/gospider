package task

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nange/gospider/common"
	"github.com/nange/gospider/spider"
	"github.com/nange/gospider/web/core"
	"github.com/nange/gospider/web/model"
	"github.com/nange/gospider/web/service"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// 根据任务id启动非定时任务
func StartTask(c *gin.Context) {
	taskIDStr := c.Param("id")
	if taskIDStr == "" {
		logrus.Warnf("StartTaskReq taskID is empty")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	taskID, err := strconv.ParseUint(taskIDStr, 10, 64)
	if err != nil {
		logrus.Warnf("StartTaskReq taskID format is invalid, taskID: %v", taskIDStr)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	logrus.Infof("StartTaskReq:%+v", taskID)

	if taskLock.IsRunning(taskID) {
		c.String(http.StatusConflict, "other operation is running")
		return
	}
	defer taskLock.Complete(taskID)

	// query task info from db
	task := &model.Task{}
	err = model.NewTaskQuerySet(core.GetDB()).IDEq(taskID).One(task)
	if err != nil {
		logrus.Errorf("StartTaskReq query model task fail, taskID: %v , err: %+v", taskIDStr, err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	// not allow crontab task
	if task.CronSpec != "" {
		logrus.Warnf("StartTaskReq taskID is crontab task, taskID: %v", taskIDStr)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// check task status
	if !taskCanBeStart(task) {
		logrus.Warnf("StartTaskReq taskID status is non-conformance , taskID: %v", taskIDStr)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// create Task Model
	spiderTask, err := service.GetSpiderTaskByModel(task)
	if err != nil {
		logrus.Errorf("StartTaskReq get model task fail, taskID: %v , err: %+v", taskIDStr, err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	// run Task Model
	err = spider.Run(spiderTask, service.GetMTSChan())
	if err != nil {
		logrus.Errorf("StartTaskReq run task fail, taskID: %v , err: %+v", taskIDStr, err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	// update task status
	err = model.NewTaskQuerySet(core.GetDB()).IDEq(taskID).GetUpdater().SetStatus(common.TaskStatusRunning).Update()
	if err != nil {
		if err := spider.CancelTask(taskID); err != nil {
			logrus.Warnf("spider.CancelTask err:%v", err)
		}
		logrus.Errorf("StartTaskReq update task status err:%+v", errors.WithStack(err))
		c.String(http.StatusInternalServerError, "")
		return
	}

	c.String(http.StatusOK, "success")
}

func taskCanBeStart(task *model.Task) bool {
	return task.Status == common.TaskStatusStopped ||
		task.Status == common.TaskStatusUnexceptedExited ||
		task.Status == common.TaskStatusCompleted ||
		task.Status == common.TaskStatusRunningTimeout
}
