package router

import (
	"github.com/gin-gonic/gin"
	"shop_api/api"
)

func CaptchaRouter(r *gin.RouterGroup) {
	app := api.App.CaptchaApi
	r.GET("base/captcha", app.CaptchaCreateView)

}
