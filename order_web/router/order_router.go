package router

import (
	"github.com/gin-gonic/gin"
	"order_web/api"
)

func OrderRouter(r *gin.RouterGroup) {
	app := api.App.OrderApi
	order := r.Group("orders")

	order.GET("", app.OrderListView)          //查看所有订单
	order.DELETE("/:id", app.DeleteOrderView) //删除订单
	order.POST("", app.OrderCreateView)       //创建订单
	//order.PATCH("/:id", app.UpdatePatchView)  // 更新订单
}
