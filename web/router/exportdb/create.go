package exportdb

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nange/gospider/web/core"
	"github.com/nange/gospider/web/model"
	log "github.com/sirupsen/logrus"
)

type CreateExportDBReq struct {
	model.ExportDB
}

type CreateExportDBResp struct {
	ID       uint64    `json:"id"`
	CreateAt time.Time `json:"create_at"`
}

func CreateExportDB(c *gin.Context) {
	var req CreateExportDBReq
	if err := c.BindJSON(&req); err != nil {
		log.Errorf("bind json failed! err:%+v", err)
		c.String(http.StatusBadRequest, "")
		return
	}
	log.Infof("req:%+v", req)

	if req.Host == "" {
		req.Host = "127.0.0.1"
	}
	if req.Port == 0 {
		req.Port = 3306
	}
	if req.User == "" {
		req.User = "root"
	}

	db := core.GetGormDB()
	if err := req.Create(db); err != nil {
		log.Errorf("create sysdb err: %+v", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, &CreateExportDBResp{
		ID:       req.ID,
		CreateAt: req.CreatedAt,
	})
}
