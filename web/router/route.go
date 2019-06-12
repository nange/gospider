package router

import (
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/gobuffalo/packr"
	"github.com/nange/gospider/web/core"
	"github.com/nange/gospider/web/model"
	"github.com/nange/gospider/web/router/exportdb"
	"github.com/nange/gospider/web/router/rule"
	"github.com/nange/gospider/web/router/task"
	"github.com/nange/gospider/web/router/user"
	log "github.com/sirupsen/logrus"
)

func Route(engine *gin.Engine) {
	authMiddleware, err := JwtAuth()
	if err != nil {
		log.Fatalf("get jwt middleware err [%+v]", err)
	}

	engine.POST("/login", authMiddleware.LoginHandler)
	api := engine.Group("/api")
	api.Use(authMiddleware.MiddlewareFunc())
	{
		api.GET("/tasks", task.GetTaskList)
		api.GET("/tasks/:id", task.GetTaskByID)
		api.POST("/tasks", task.CreateTask)
		api.PUT("/tasks/:id", task.UpdateTask)
		api.PUT("/tasks/:id/stop", task.StopTask)
		api.PUT("/tasks/:id/start", task.StartTask)
		api.PUT("/tasks/:id/restart", task.RestartTask)

		api.GET("/rules", rule.GetRuleList)

		api.GET("/exportdb", exportdb.GetExportDBList)
		api.POST("/exportdb", exportdb.CreateExportDB)
		api.DELETE("/exportdb/:id", exportdb.DeleteExportDB)

		api.GET("/user/info", user.GetUserInfo)

	}

	box := packr.NewBox("../static/dist")
	box2 := packr.NewBox("../static/dist/static")
	engine.StaticFS("/admin", box)
	engine.StaticFS("/static", box2)
}

func JwtAuth() (*jwt.GinJWTMiddleware, error) {
	type login struct {
		Username string `form:"username" json:"username" binding:"required"`
		Password string `form:"password" json:"password" binding:"required"`
	}
	middle, err := jwt.New(&jwt.GinJWTMiddleware{
		Key:        []byte("gospider"),
		Timeout:    24 * time.Hour,
		MaxRefresh: 24 * time.Hour,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*model.User); ok {
				return jwt.MapClaims{
					jwt.IdentityKey: v,
				}
			}
			return jwt.MapClaims{}
		},
		Authenticator: func(c *gin.Context) (i interface{}, e error) {
			var loginVals login
			if err := c.ShouldBind(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			valid, user, err := model.IsValidUser(core.GetGormDB(), loginVals.Username, loginVals.Password)
			if err != nil {
				return nil, err
			}
			if valid {
				return user, nil
			}

			return nil, jwt.ErrFailedAuthentication
		},
		TokenLookup: "header: Authorization, query: token, cookie: jwt",
		SendCookie:  true,
	})

	return middle, err
}
