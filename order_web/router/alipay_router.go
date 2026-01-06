package router

import (
	"github.com/gin-gonic/gin"
	"order_web/api"
	"order_web/middleware"
)

func AlipayRouter(r *gin.RouterGroup) {
	app := api.App.OrderApi

	alipay := r.Group("/pay").Use(middleware.Trace)
	alipay.POST("/callback", app.AlipayCallBackView)

}
