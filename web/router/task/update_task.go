package task

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nange/gospider/common"
	"github.com/nange/gospider/web/core"
	"github.com/nange/gospider/web/model"
	"github.com/nange/gospider/web/service"
	"github.com/sirupsen/logrus"
)

type UpdateTaskReq struct {
	model.Task
	OutputSysDBID string `json:"sysdb_id"`
}

type UpdateTaskResp struct {
	ID       uint64    `json:"id"`
	UpdateAt time.Time `json:"update_at"`
}

func UpdateTask(c *gin.Context) {
	taskIDStr := c.Param("id")
	if taskIDStr == "" {
		logrus.Warnf("UpdateTaskReq taskID is empty")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	taskID, err := strconv.ParseUint(taskIDStr, 10, 64)
	if err != nil {
		logrus.Warnf("UpdateTaskReq taskID format is invalid, taskID: %v", taskIDStr)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var req UpdateTaskReq
	if err := c.BindJSON(&req); err != nil {
		logrus.Errorf("UpdateTaskReq bind json failed! err:%+v", err)
		c.String(http.StatusBadRequest, "")
		return
	}
	logrus.Infof("UpdateTaskreq:%+v %+v", taskID, req)
	req.Task.ID = taskID

	task := req.Task

	if taskLock.IsRunning(taskID) {
		c.String(http.StatusConflict, "other operation is running")
		return
	}
	defer taskLock.Complete(taskID)

	//get current task infp
	oldtask := &model.Task{}
	err = model.NewTaskQuerySet(core.GetDB()).IDEq(taskID).One(oldtask)
	if err != nil {
		logrus.Errorf("UpdateTaskReq query model task fail, taskID: %v , err: %+v", taskID, err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// check task status
	if !taskCanBeUpdate(oldtask) {
		logrus.Warnf("UpdateTaskReq taskID status is non-conformance , taskID: %v", taskID)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	// recreate crontab task if cronspec is change
	if err := cronTaskStopAndCreate(taskID, oldtask.Status, *oldtask, task); err != nil {
		logrus.Errorf("UpdateTaskReq recreate crontab task fail, taskID: %v , err: %+v", taskID, err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	// update db
	if err := task.Update(core.GetDB(),
		model.TaskDBSchema.TaskDesc, model.TaskDBSchema.CronSpec,
		model.TaskDBSchema.OutputType, model.TaskDBSchema.OutputSysDBID,
		model.TaskDBSchema.OptUserAgent, model.TaskDBSchema.OptMaxDepth,
		model.TaskDBSchema.OptAllowedDomains, model.TaskDBSchema.OptURLFilters,
		model.TaskDBSchema.OptMaxBodySize, model.TaskDBSchema.LimitEnable,
		model.TaskDBSchema.LimitDomainRegexp, model.TaskDBSchema.LimitDomainGlob,
		model.TaskDBSchema.LimitDelay, model.TaskDBSchema.LimitRandomDelay,
		model.TaskDBSchema.LimitParallelism, model.TaskDBSchema.ProxyURLs,
		model.TaskDBSchema.OptRequestTimeout,
	); err != nil {
		// task roll back
		if err := cronTaskStopAndCreate(taskID, oldtask.Status, task, *oldtask); err != nil {
			logrus.Errorf("UpdateTaskReq rollback crontab task fail, taskID: %v , err: %+v", taskID, err)
		}
		logrus.Errorf("UpdateTaskReq update task failed! err:%+v", err)
		c.String(http.StatusInternalServerError, "")
		return
	}

	c.JSON(http.StatusOK, &UpdateTaskResp{
		ID:       task.ID,
		UpdateAt: task.UpdatedAt,
	})
}

func taskCanBeUpdate(task *model.Task) bool {
	return task.Status == common.TaskStatusStopped ||
		task.Status == common.TaskStatusUnexceptedExited ||
		task.Status == common.TaskStatusCompleted ||
		task.Status == common.TaskStatusRunningTimeout ||
		task.Status == common.TaskStatusRunning

}

func cronTaskStopAndCreate(taskID uint64, taskStatus common.TaskStatus, oldtask, newtask model.Task) error {
	if oldtask.CronSpec == newtask.CronSpec {
		return nil
	}
	if oldtask.CronSpec != "" {
		// stop cron task
		if ct := service.GetCronTask(taskID); ct != nil {
			ct.Stop()
		}
	}
	if newtask.CronSpec != "" && (taskStatus == common.TaskStatusCompleted || taskStatus == common.TaskStatusUnexceptedExited) {
		return service.CreateCronTask(newtask)
	}
	return nil
}
