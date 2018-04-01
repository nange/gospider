package task

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nange/gospider/web/common"
	"github.com/nange/gospider/web/model"
	"github.com/sirupsen/logrus"
)

type GetTaskListResp struct {
	Total int          `json:"total"`
	Data  []model.Task `json:"data"`
}

func GetTaskList(c *gin.Context) {
	db, err := common.GetGormDBFromEnv()
	if err != nil {
		logrus.Errorf("GetDBFromEnv failed! err:%#v", err)
		c.Data(http.StatusInternalServerError, "", nil)
		return
	}

	tasks, count, err := model.GetTaskList(db, -1, -1)
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
