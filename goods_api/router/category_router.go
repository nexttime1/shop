package router

import (
	"github.com/gin-gonic/gin"
	"goods_api/api"
)

func CategoryRouter(r *gin.RouterGroup) {
	app := api.App.CategoryApi
	r.GET("categorys", app.GetAllCategoryView)
	r.GET("categorys/:id", app.GetSubCategoryView)
	r.POST("categorys", app.CreateCategoryView)
	r.PUT("categorys/:id", app.UpdateCategoryView)
	r.DELETE("categorys/:id", app.DeleteCategoryView)
}
