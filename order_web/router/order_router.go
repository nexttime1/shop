package router

import (
	"github.com/gin-gonic/gin"
	"order_web/api"
	"order_web/middleware"
)

func OrderRouter(r *gin.RouterGroup) {
	app := api.App.OrderApi
	order := r.Group("orders").Use(middleware.AuthMiddleware).Use(middleware.Trace)
	// 限流
	order.GET("", middleware.OrderListCurrentLimiting, app.OrderListView)      //查看所有订单
	order.POST("", middleware.CreateOrderCurrentLimiting, app.OrderCreateView) //创建订单
	order.GET("/:id", app.OrderDetailView)                                     //订单细节

}
