package router

import (
	"github.com/gin-gonic/gin"
	"shop_api/api"
)

func CaptchaRouter(r *gin.RouterGroup) {
	app := api.App.CaptchaApi
	r.GET("base/captcha", app.CaptchaCreateView)
	r.POST("base/send_sms", app.SendRegisterView)
	r.POST("base/verify_sms", app.VerifyCaptchaView) //测试用

}
