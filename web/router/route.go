package router

import (
	"github.com/gin-gonic/gin"
	"github.com/gobuffalo/packr"
	"github.com/nange/gospider/web/router/rule"
	"github.com/nange/gospider/web/router/sysdb"
	"github.com/nange/gospider/web/router/task"
)

func Route(engine *gin.Engine) {
	api := engine.Group("/api")
	{
		api.GET("/tasks", task.GetTaskList)
		api.GET("/tasks/:id", task.GetTaskByID)
		api.POST("/tasks", task.CreateTask)
		api.PUT("/tasks/:id", task.UpdateTask)
		api.PUT("/tasks/:id/stop", task.StopTask)
		api.PUT("/tasks/:id/start", task.StartTask)
		api.PUT("/tasks/:id/restart", task.RestartTask)

		api.GET("/rules", rule.GetRuleList)

		api.GET("/sysdbs", sysdb.GetSysDBs)
		api.POST("/sysdbs", sysdb.CreateSysDB)

	}

	box := packr.NewBox("../static/dist")
	box2 := packr.NewBox("../static/dist/static")
	engine.StaticFS("/admin", box)
	engine.StaticFS("/static", box2)
}
