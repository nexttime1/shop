package router

import (
	"github.com/gin-gonic/gin"
	"order_web/api"
)

func AlipayRouter(r *gin.RouterGroup) {
	app := api.App.OrderApi
	alipay := r.Group("/pay")
	alipay.POST("/callback", app.AlipayCallBackView)

}
