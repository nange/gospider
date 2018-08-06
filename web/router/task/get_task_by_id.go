package task

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nange/gospider/web/core"
	"github.com/nange/gospider/web/model"
	"github.com/sirupsen/logrus"
)

type GetTaskIDReq struct {
	ID uint64 `json:"id" in:"path" validate:"@uint64[1,]"`
}

type GetTaskIDResp struct {
	Total int        `json:"total"`
	Data  model.Task `json:"data"`
}

func GetTaskByID(c *gin.Context) {
	taskIDStr := c.Param("id")
	if taskIDStr == "" {
		logrus.Warnf("GetTaskByID taskID is empty")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	taskID, err := strconv.ParseUint(taskIDStr, 10, 64)
	if err != nil {
		logrus.Warnf("GetTaskByID taskID format is invalid, taskID:%v", taskIDStr)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	logrus.Infof("GetTaskByID:%+v", taskID)

	// query task info from db
	task := &model.Task{}
	err = model.NewTaskQuerySet(core.GetDB()).IDEq(taskID).One(task)
	if err != nil {
		logrus.Errorf("GetTaskByID query model task fail, taskID: %v , err: %+v", taskIDStr, err)
		c.String(http.StatusInternalServerError, "")
		return
	}
	resp := &GetTaskIDResp{
		Total: 1,
		Data:  *task,
	}
	c.JSON(http.StatusOK, resp)
}
