package router

import (
	"github.com/gin-gonic/gin"
	"user_web/api"
	"user_web/middleware"
)

func CaptchaRouter(r *gin.RouterGroup) {
	app := api.App.CaptchaApi
	r.Use(middleware.Trace)
	r.GET("base/captcha", app.CaptchaCreateView)
	r.POST("base/send_sms", app.SendRegisterView)
	r.POST("base/verify_sms", app.VerifyCaptchaView) //测试用

}
