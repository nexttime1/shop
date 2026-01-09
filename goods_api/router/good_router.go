package router

import (
	"github.com/gin-gonic/gin"
	"goods_api/api"
	"goods_api/middleware"
)

func GoodRouter(r *gin.RouterGroup) {
	app := api.App.GoodApi
	r.Use(middleware.Trace)                                             //链路追踪
	r.GET("good/list", middleware.CurrentLimiting, app.GetGoodListView) // 限流
	r.POST("good", app.CreateGoodView)
	r.GET("good/:id", app.GoodDetailView)
	r.PUT("good/:id", app.GoodUpdateView)
	r.PATCH("good/:id", app.GoodPatchUpdateView)
	r.DELETE("good/:id", app.GoodDeleteView)
}
