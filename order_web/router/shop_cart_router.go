package router

import (
	"github.com/gin-gonic/gin"
	"order_web/api"
	"order_web/middleware"
)

func CartRouter(r *gin.RouterGroup) {
	app := api.App.CartApi
	cart := r.Group("shopcarts").Use(middleware.AuthMiddleware)

	cart.GET("", app.CartListView)              //购物车列表
	cart.DELETE("/:id", app.DeleteCartItemView) //删除条目
	cart.POST("", app.AddItemView)              //添加商品到购物车
	cart.PATCH("/:id", app.UpdatePatchView)     // 更新购物车中的某个商品
}
