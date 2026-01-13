package router

import (
	"github.com/gin-gonic/gin"
	"goods_web/api"
	"goods_web/middleware"
)

func PmsRouter(r *gin.Engine) {
	app := api.App.PmsApi
	g := r.Group("/pms/v1").Use(middleware.Trace)
	g.GET("/productAttr", app.ProductAttrListView)
	g.POST("/productAttr", app.ProductAttrCreateView)
	g.GET("/skuStock", app.SkuStockListView)
	g.POST("/skuStock", app.SkuStockCreateView)
}
