package router

import (
	"github.com/gin-gonic/gin"
	"github.com/gobuffalo/packr"
	"github.com/nange/gospider/web/router/rule"
	"github.com/nange/gospider/web/router/task"
)

func Route(engine *gin.Engine) {
	api := engine.Group("/api")
	{
		api.GET("/tasks", task.GetTaskList)
		api.POST("/tasks", task.CreateTask)

		api.GET("/rules", rule.GetRuleList)
	}

	box := packr.NewBox("../static/dist")
	box2 := packr.NewBox("../static/dist/static")
	engine.StaticFS("/admin", box)
	engine.StaticFS("/static", box2)
}
