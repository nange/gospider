package web

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/nange/gospider/web/model"
	"github.com/nange/gospider/web/router"
	"github.com/pkg/errors"
)

type Server struct {
	IP   string
	Port int
	db   *gorm.DB
}

func (s *Server) SetDB(gdb *gorm.DB) {
	s.db = gdb
}

func (s *Server) Run() error {
	model.SetDB(s.db)
	if err := model.AutoMigrate(); err != nil {
		return errors.Wrap(err, "model auto migrate failed")
	}

	engine := gin.Default()
	router.Route(engine)

	return errors.WithStack(engine.Run(fmt.Sprintf("%s:%d", s.IP, s.Port)))
}
