package router

import (
	"github.com/gin-gonic/gin"
	"option_web/api"
	"option_web/middleware"
)

func MessageRouter(r *gin.RouterGroup) {
	app := api.App.MessageApi
	message := r.Group("message").Use(middleware.AuthMiddleware).Use(middleware.Trace) //跟踪
	message.GET("", app.MessageListView)                                               // 消息列表
	message.POST("", app.CreateMessageView)                                            //添加留言
}
