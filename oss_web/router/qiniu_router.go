package router

import (
	"github.com/gin-gonic/gin"
	"oss_web/common/res"
	"oss_web/global"
	"oss_web/middleware"
	"oss_web/utils/qiniu"
)

func QiNiuRouter(r *gin.RouterGroup) {
	r.Use(middleware.Trace)
	r.GET("/token", func(c *gin.Context) {
		filename := c.Query("filename")
		if filename == "" {
			res.FailWithMsg(c, res.FailArgumentCode, "filename不能为空")
			return
		}
		tokenData, err := qiniu.GetUploadTokenForBrowser(global.Config.QiNiu.Prefix, filename)
		if err != nil {
			res.FailWithMsg(c, res.FailServiceCode, err.Error())
			return
		}
		res.OkWithData(c, tokenData)
	})
	// 七牛云回调接口
	r.POST("/callback", qiniu.QiNiuCallback)
}
