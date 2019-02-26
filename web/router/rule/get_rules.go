package rule

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nange/gospider/spider"
	log "github.com/sirupsen/logrus"
)

type GetRuleListResp struct {
	Total int      `json:"total"`
	Data  []string `json:"data"`
}

func GetRuleList(c *gin.Context) {
	keys := spider.GetTaskRuleKeys()
	if len(keys) == 0 {
		log.Warnf("task rule is empty")
	}

	c.JSON(http.StatusOK, &GetRuleListResp{
		Total: len(keys),
		Data:  keys,
	})
}
