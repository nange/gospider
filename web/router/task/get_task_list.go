package task

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nange/gospider/web/model"
	"github.com/sirupsen/logrus"
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
		logrus.Warnf("query param is invalid")
		c.Data(http.StatusBadRequest, "", nil)
		return
	}
	logrus.Infof("get task list req:%+v", req)

	tasks, count, err := model.GetTaskList(req.Size, req.Offset)
	if err != nil {
		logrus.Errorf("GetTaskList failed! err:%#v", err)
		c.Data(http.StatusInternalServerError, "", nil)
		return
	}

	c.JSON(http.StatusOK, &GetTaskListResp{
		Total: count,
		Data:  tasks,
	})
}
