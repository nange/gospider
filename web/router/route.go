package router

import (
	"github.com/gin-gonic/gin"
	"github.com/nange/gospider/web/router/task"
)

func Route(engine *gin.Engine)  {
	api := engine.Group("/api")
	{
		api.GET("/tasks", task.GetTaskList)
	}
}
