package router

import (
	"github.com/gin-gonic/gin"
	"option_web/api"
	"option_web/middleware"
)

func CollectionRouter(r *gin.RouterGroup) {
	app := api.App.CollectionApi
	userfavs := r.Group("userfavs").Use(middleware.AuthMiddleware).Use(middleware.Trace) //跟踪
	userfavs.GET("", app.CollectionListView)                                             //查看收藏
	userfavs.DELETE("/:good_id", app.CollectionDeleteView)                               //删除收藏
	userfavs.POST("", app.CollectionAddView)                                             //添加收藏
	userfavs.GET("/:good_id", app.CollectionDetailView)                                  //查看详情

}
