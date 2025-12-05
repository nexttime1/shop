package router

import (
	"github.com/gin-gonic/gin"
	"goods_api/api"
)

func GoodRouter(r *gin.RouterGroup) {
	app := api.App.GoodApi
	r.GET("good/list", app.GetGoodList)

}
