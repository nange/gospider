package sysdb

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nange/gospider/web/core"
	"github.com/nange/gospider/web/model"
	"github.com/sirupsen/logrus"
)

type CreateSysDBReq struct {
	model.SysDB
}

type CreateSysDBResp struct {
	ID       uint64    `json:"id"`
	CreateAt time.Time `json:"create_at"`
}

func CreateSysDB(c *gin.Context) {
	var req CreateSysDBReq
	if err := c.BindJSON(&req); err != nil {
		logrus.Errorf("bind json failed! err:%+v", err)
		c.Data(http.StatusBadRequest, "", nil)
		return
	}
	logrus.Infof("req:%+v", req)

	db := core.GetDB()
	if err := req.Create(db); err != nil {
		logrus.Errorf("create sysdb err: %+v", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, &CreateSysDBResp{
		ID:       req.ID,
		CreateAt: req.CreatedAt,
	})
}
