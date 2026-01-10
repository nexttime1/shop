package router

import (
	"github.com/gin-gonic/gin"
	"user_web/api"
	"user_web/middleware"
)

func UserRouter(r *gin.RouterGroup) {
	app := api.App.UserApi
	r.Use(middleware.Trace)
	r.GET("/user/list", middleware.AdminMiddleware, app.UserListView)
	r.POST("/user/login", app.UserLoginView)
	r.POST("/user/register", app.UserRegisterView)
}
