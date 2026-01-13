package router

import (
	"github.com/gin-gonic/gin"
	"order_web/api"
	"order_web/middleware"
)

func SmsRouter(r *gin.Engine) {
	app := api.App.SmsApi
	g := r.Group("/sms/v1").Use(middleware.Trace)
	// 优惠券模块
	g.GET("/coupons", app.CouponListView)
	g.GET("/coupons/:id", app.CouponDetailView)
	g.POST("/coupons", app.CouponCreateView)
	// 秒杀活动模块
	g.GET("/flash", app.FlashListView)
	g.GET("/flash/:id", app.FlashDetailView)
	g.POST("/flash", app.FlashCreateView)
	// 广告位模块
	g.GET("/ads", app.AdListView)
	g.GET("/ads/:id", app.AdDetailView)
	g.POST("/ads", app.AdCreateView)
}
