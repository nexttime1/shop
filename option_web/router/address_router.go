package router

import (
	"github.com/gin-gonic/gin"
	"option_web/api"
	"option_web/middleware"
)

func AddressRouter(r *gin.RouterGroup) {
	app := api.App.AddressApi
	order := r.Group("address").Use(middleware.AuthMiddleware)

	order.GET("", app.AddressListView)          //查看所有地址
	order.DELETE("/:id", app.DeleteAddressView) //删除地址
	order.POST("", app.AddressCreateView)       //创建地址
	order.PUT("/:id", app.UpdateAddressView)    //修改地址

}
