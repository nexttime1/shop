package router

import (
	"github.com/gin-gonic/gin"
	"goods_web/api"
	"goods_web/middleware"
)

func BrandRouter(r *gin.RouterGroup) {
	var app = api.App.BrandApi
	r.Use(middleware.Trace) //链路追踪
	// 品牌自身
	r.GET("brands", app.BrandListView)
	r.POST("brands", app.CreateBrandView)
	r.PUT("brands/:id", app.UpdateBrandView)
	r.DELETE("brands/:id", app.DeleteBrandView)
	r.GET("brands/:id", app.BrandDetailView)

	// 第三张表
	r.GET("categorybrands", app.CategoryBrandListView)    //所有的 第三张表
	r.GET("categorybrands/:id", app.CategoryAllBrandView) //某个分类下的所有品牌
	r.POST("categorybrands", app.CreateCategoryBrandView)
	r.PUT("categorybrands/:id", app.UpdateCategoryBrandView)
	r.DELETE("categorybrands/:id", app.DeleteCategoryBrandView)

}
