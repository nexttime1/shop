package router

import (
	"github.com/gin-gonic/gin"
	"goods_api/api"
)

func BannerRouter(r *gin.RouterGroup) {
	app := api.App.BannerApi
	r.GET("banners", app.GetBannerListView)
	r.POST("banners", app.CreateBannerView)
	r.PUT("banners/:id", app.UpdateBannerView)
	r.DELETE("banners/:id", app.DeleteBannerView)
}
