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
	log "github.com/sirupsen/logrus"
)

// 根据任务id停止任务
func StopTask(c *gin.Context) {
	taskIDStr := c.Param("id")
	if taskIDStr == "" {
		log.Warnf("StopTaskReq taskID is empty")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	taskID, err := strconv.ParseUint(taskIDStr, 10, 64)
	if err != nil {
		log.Warnf("StopTaskReq taskID format is invalid, taskID:%v", taskIDStr)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	log.Infof("StopTaskReq:%+v", taskID)

	if taskLock.IsRunning(taskID) {
		c.String(http.StatusConflict, "other operation is running")
		return
	}
	defer taskLock.Complete(taskID)

	// stop cron task
	if ct := service.GetCronTask(taskID); ct != nil {
		ct.Stop()
	}

	// cancel spider task
	spider.CancelTask(taskID)

	// set task status to TaskStatusStopped
	err = model.NewTaskQuerySet(core.GetGormDB()).IDEq(taskID).GetUpdater().SetStatus(common.TaskStatusStopped).Update()
	if err != nil {
		log.Errorf("update task status err:%+v", errors.WithStack(err))
		c.String(http.StatusInternalServerError, "")
		return
	}

	c.String(http.StatusOK, "success")
}
