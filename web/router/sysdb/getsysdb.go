package sysdb

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nange/gospider/web/model"
	"github.com/sirupsen/logrus"
)

type GetSysDBListReq struct {
	Size   int `json:"size" form:"size"`
	Offset int `json:"offset" form:"offset"`
}

type GetSysDBListResp struct {
	Total int           `json:"total"`
	Data  []model.SysDB `json:"data"`
}

func GetSysDBs(c *gin.Context) {
	var req GetSysDBListReq
	if err := c.BindQuery(&req); err != nil {
		logrus.Warnf("query param is invalid")
		c.String(http.StatusBadRequest, "")
		return
	}
	logrus.Infof("get sysdb list req:%+v", req)

	tasks, count, err := model.GetSysDBList(req.Size, req.Offset)
	if err != nil {
		logrus.Errorf("GetSysDBList failed! err:%+v", err)
		c.String(http.StatusInternalServerError, "")
		return
	}

	c.JSON(http.StatusOK, &GetSysDBListResp{
		Total: count,
		Data:  tasks,
	})
}
