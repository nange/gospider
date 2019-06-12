package user

import (
	"net/http"
	"strings"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func GetUserInfo(c *gin.Context) {
	v, exists := c.Get(jwt.IdentityKey)
	if !exists {
		log.Warnf("jwt identity key not found")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	m := v.(map[string]interface{})
	log.Debugf("identity:%+v", m)

	rolesStr := m["roles"].(string)
	m["roles"] = strings.Split(rolesStr, ",")

	c.JSON(http.StatusOK, m)
}
