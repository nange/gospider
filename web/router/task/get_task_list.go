package task

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nange/gospider/web/core"
	"github.com/nange/gospider/web/model"
	log "github.com/sirupsen/logrus"
)

type GetTaskListReq struct {
	Size   int `json:"size" form:"size"`
	Offset int `json:"offset" form:"offset"`
}

type GetTaskListResp struct {
	Total int          `json:"total"`
	Data  []model.Task `json:"data"`
}

func GetTaskList(c *gin.Context) {
	var req GetTaskListReq
	if err := c.BindQuery(&req); err != nil {
		log.Warnf("query param is invalid")
		c.String(http.StatusBadRequest, "")
		return
	}
	log.Infof("get task list req:%+v", req)

	tasks, count, err := model.GetTaskList(core.GetGormDB(), req.Size, req.Offset)
	if err != nil {
		log.Errorf("GetTaskList failed! err:%+v", err)
		c.String(http.StatusInternalServerError, "")
		return
	}

	c.JSON(http.StatusOK, &GetTaskListResp{
		Total: count,
		Data:  tasks,
	})
}
