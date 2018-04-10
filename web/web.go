package web

import (
	"github.com/gin-gonic/gin"
	"github.com/nange/gospider/web/common"
	"github.com/nange/gospider/web/model"
	"github.com/nange/gospider/web/router"
	"github.com/sirupsen/logrus"
)

func Run() {
	model.AutoMigrate(common.GetDB())

	engine := gin.Default()
	router.Route(engine)

	logrus.Fatal(engine.Run())
}
