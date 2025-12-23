package router

import (
	"github.com/gin-gonic/gin"
	"option_web/api"
	"option_web/middleware"
)

func CollectionRouter(r *gin.RouterGroup) {
	app := api.App.Collection
	userfavs := r.Group("userfavs").Use(middleware.AuthMiddleware)
	userfavs.GET("", app.AddressListView)          //查看所有地址
	userfavs.DELETE("/:id", app.DeleteAddressView) //删除地址
	userfavs.POST("", app.AddressCreateView)       //创建地址
	userfavs.GET("/:id", app.AddressListView)      //查看所有地址

}
