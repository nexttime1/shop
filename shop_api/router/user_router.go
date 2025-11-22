package router

import (
	"github.com/gin-gonic/gin"
	"shop_api/api"
)

func UserRouter(r *gin.RouterGroup) {
	app := api.App.UserApi
	r.GET("/user/list", app.UserListView)
	r.POST("/user/login", app.UserLoginView)
}
