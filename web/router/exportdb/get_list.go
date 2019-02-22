package exportdb

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nange/gospider/web/core"
	"github.com/nange/gospider/web/model"
	log "github.com/sirupsen/logrus"
)

type GetExportDBListReq struct {
	Size   int `json:"size" form:"size"`
	Offset int `json:"offset" form:"offset"`
}

type GetExportDBListResp struct {
	Total int              `json:"total"`
	Data  []model.ExportDB `json:"data"`
}

func GetExportDBList(c *gin.Context) {
	var req GetExportDBListReq
	if err := c.BindQuery(&req); err != nil {
		log.Warnf("query param is invalid")
		c.String(http.StatusBadRequest, "")
		return
	}
	log.Infof("get sysdb list req:%+v", req)

	tasks, count, err := model.GetExportDBList(core.GetDB(), req.Size, req.Offset)
	if err != nil {
		log.Errorf("GetExportDBList failed! err [%+v]", err)
		c.String(http.StatusInternalServerError, "")
		return
	}

	c.JSON(http.StatusOK, &GetExportDBListResp{
		Total: count,
		Data:  tasks,
	})
}
