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
	var req UpdateTaskReq
	if err := c.BindJSON(&req); err != nil {
		logrus.Errorf("UpdateTaskReq bind json failed! err:%+v", err)
		c.String(http.StatusBadRequest, "")
		return
	}
	logrus.Infof("req:%+v", req)

	intID, err := strconv.Atoi(req.OutputSysDBID)
	if err != nil {
		c.String(http.StatusBadRequest, "")
		return
	}
	req.Task.OutputSysDBID = uint64(intID)
	task := req.Task

	taskID := task.ID
	if taskLock.IsRunning(taskID) {
		c.String(http.StatusConflict, "任务正在执行")
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
	if err := task.Update(core.GetDB(), model.TaskDBSchemaField("id"),
		model.TaskDBSchemaField("task_desc"), model.TaskDBSchemaField("cron_spec"),
		model.TaskDBSchemaField("output_type"), model.TaskDBSchemaField("output_sysdb_id"),
		model.TaskDBSchemaField("opt_user_agent"), model.TaskDBSchemaField("opt_max_depth"),
		model.TaskDBSchemaField("opt_allowed_domains"), model.TaskDBSchemaField("opt_url_filters"),
		model.TaskDBSchemaField("opt_max_body_size"), model.TaskDBSchemaField("limit_enable"),
		model.TaskDBSchemaField("limit_domain_regexp"), model.TaskDBSchemaField("limit_domain_glob"),
		model.TaskDBSchemaField("limit_delay"), model.TaskDBSchemaField("limit_random_delay"),
		model.TaskDBSchemaField("limit_parallelism"), model.TaskDBSchemaField("proxy_urls"),
	); err != nil {
		// task roll back
		if err := cronTaskStopAndCreate(taskID, oldtask.Status, task, *oldtask); err != nil {
			logrus.Errorf("UpdateTaskReq rollback crontab task fail, taskID: %v , err: %+v", taskID, err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
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
		task.Status == common.TaskStatusRunningTimeout
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
