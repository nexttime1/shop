package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"shop_api/global"
)

func Trace(c *gin.Context) {

	parentScan := global.Tracer.StartSpan(c.Request.URL.Path) //用路径当作 开始
	zap.L().Info(c.Request.URL.Path)
	defer parentScan.Finish()
	c.Set("parent_scan", parentScan)
	zap.S().Infof("parent_scan: %v", parentScan)
	zap.S().Infof("ServiceName 为 %s", global.Config.JaegerInfo.ServiceName)
	c.Next()

}
