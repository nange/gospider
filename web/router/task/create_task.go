package task

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nange/gospider/common"
	"github.com/nange/gospider/spider"
	"github.com/nange/gospider/web/core"
	"github.com/nange/gospider/web/model"
	"github.com/nange/gospider/web/service"
	"github.com/sirupsen/logrus"
)

type CreateTaskReq struct {
	model.Task
	OutputSysDBID string `json:"sysdb_id"`
}

type CreateTaskResp struct {
	ID       uint64    `json:"id"`
	CreateAt time.Time `json:"create_at"`
}

// TODO: 定时任务停止功能
func CreateTask(c *gin.Context) {
	var req CreateTaskReq
	if err := c.BindJSON(&req); err != nil {
		logrus.Errorf("bind json failed! err:%+v", err)
		c.Data(http.StatusBadRequest, "", nil)
		return
	}
	logrus.Infof("req:%+v", req)

	intID, err := strconv.Atoi(req.OutputSysDBID)
	if err != nil {
		c.Data(http.StatusBadRequest, "", nil)
		return
	}
	req.Task.OutputSysDBID = uint64(intID)

	task := req.Task
	task.Status = common.TaskStatusRunning
	if err := task.Create(core.GetDB()); err != nil {
		logrus.Errorf("create task failed! err:%+v", err)
		c.Data(http.StatusInternalServerError, "", nil)
		return
	}

	spiderTask, err := service.GetSpiderTaskByModel(&task)
	err = spider.Run(spiderTask, service.GetMTSChan())
	if err != nil {
		logrus.Errorf("spider run task failed! err:%+v", err)
		c.Data(http.StatusInternalServerError, "", nil)
		return
	}

	if task.CronSpec != "" {
		logrus.Infof("starting cron task:%s", task.CronSpec)
		ct, err := service.NewCronTask(&task, service.GetMTSChan())
		if err != nil {
			logrus.Errorf("new cron task failed! err:%+v", err)
		} else {
			if err := ct.Start(); err != nil {
				logrus.Errorf("start cron task failed! err:%+v", err)
			}
		}

	}

	c.JSON(http.StatusOK, &CreateTaskResp{
		ID:       task.ID,
		CreateAt: task.CreatedAt,
	})
}
