package router

import (
	"github.com/gin-gonic/gin"
	"shop_api/api"
	"shop_api/middleware"
)

func UserRouter(r *gin.RouterGroup) {
	app := api.App.UserApi
	r.GET("/user/list", middleware.AdminMiddleware, app.UserListView)
	r.POST("/user/login", app.UserLoginView)
	r.POST("/user/register", app.UserRegisterView)
}
