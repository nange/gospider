package task

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nange/gospider/common"
	"github.com/nange/gospider/spider"
	"github.com/nange/gospider/web/core"
	"github.com/nange/gospider/web/model"
	"github.com/nange/gospider/web/service"
	log "github.com/sirupsen/logrus"
)

type CreateTaskReq struct {
	model.Task
}

type CreateTaskResp struct {
	ID       uint64    `json:"id"`
	CreateAt time.Time `json:"create_at"`
}

func CreateTask(c *gin.Context) {
	var req CreateTaskReq
	if err := c.BindJSON(&req); err != nil {
		log.Errorf("bind json failed! err:%+v", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	log.Infof("req:%+v", req)

	task := req.Task
	task.Status = common.TaskStatusRunning
	if err := task.Create(core.GetGormDB()); err != nil {
		log.Errorf("create task failed! err:%+v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	spiderTask, err := service.GetSpiderTaskByModel(&task)
	if err != nil {
		log.Errorf("spider get model task failed! err:%+v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	s := spider.New(spiderTask, service.GetMTSChan())
	if err := s.Run(); err != nil {
		log.Errorf("spider run task failed! err:%+v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if task.CronSpec != "" {
		log.Infof("starting cron task:%s", task.CronSpec)
		ct, err := service.NewCronTask(task.ID, task.CronSpec, service.GetMTSChan())
		if err != nil {
			log.Errorf("new cron task failed! err:%+v", err)
		} else {
			if err := ct.Start(); err != nil {
				log.Errorf("start cron task failed! err:%+v", err)
			}
		}

	}

	c.JSON(http.StatusOK, &CreateTaskResp{
		ID:       task.ID,
		CreateAt: task.CreatedAt,
	})
}
