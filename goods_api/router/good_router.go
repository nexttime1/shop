package router

import (
	"github.com/gin-gonic/gin"
	"goods_api/api"
)

func GoodRouter(r *gin.RouterGroup) {
	app := api.App.GoodApi
	r.GET("good/list", app.GetGoodListView)
	r.POST("good", app.CreateGoodView)
	r.GET("good/:id", app.GoodDetailView)
	r.PUT("good/:id", app.GoodUpdateView)
	r.PATCH("good/:id", app.GoodPatchUpdateView)
	r.DELETE("good/:id", app.GoodDeleteView)
}
