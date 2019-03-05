package task

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nange/gospider/web/core"
	"github.com/nange/gospider/web/model"
	"github.com/sirupsen/logrus"
)

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
	err = model.NewTaskQuerySet(core.GetGormDB()).IDEq(taskID).One(task)
	if err != nil {
		logrus.Errorf("GetTaskByID query model task fail, taskID: %v , err: %+v", taskIDStr, err)
		c.String(http.StatusInternalServerError, "")
		return
	}

	c.JSON(http.StatusOK, task)
}
