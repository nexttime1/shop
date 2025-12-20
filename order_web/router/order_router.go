package router

import (
	"github.com/gin-gonic/gin"
	"order_web/api"
	"order_web/middleware"
)

func OrderRouter(r *gin.RouterGroup) {
	app := api.App.OrderApi
	order := r.Group("orders").Use(middleware.AuthMiddleware)

	order.GET("", app.OrderListView)          //查看所有订单
	order.DELETE("/:id", app.DeleteOrderView) //删除订单
	order.POST("", app.OrderCreateView)       //创建订单
	order.GET("/:id", app.OrderDetailView)    //订单细节

}
