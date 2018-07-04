package web

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/nange/gospider/web/core"
	"github.com/nange/gospider/web/router"
	"github.com/nange/gospider/web/service"
	"github.com/pkg/errors"
)

type Server struct {
	IP   string
	Port int
	db   *gorm.DB
}

func (s *Server) SetDB(gdb *gorm.DB) {
	s.db = gdb
	s.db.LogMode(true)
}

func (s *Server) Run() error {
	core.SetDB(s.db)
	if err := core.AutoMigrate(); err != nil {
		return errors.Wrap(err, "model auto migrate failed")
	}

	// 启动服务时，先检查task相关状态
	go service.CheckTask()
	// 管理task状态(如task运行完成之后需要将状态标为已完成)
	go service.ManageTaskStatus()

	engine := gin.Default()
	router.Route(engine)

	return errors.WithStack(engine.Run(fmt.Sprintf("%s:%d", s.IP, s.Port)))
}
