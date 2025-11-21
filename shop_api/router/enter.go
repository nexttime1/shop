package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"shop_api/global"
)

func Router() {
	gin.SetMode(global.Config.System.GinMode)
	r := gin.Default()
	userGroup := r.Group("/u/v1")
	UserRouter(userGroup)
	zap.L().Info("router is running ...")
	err := r.Run(global.Config.System.GetAddr())
	if err != nil {
		zap.L().Error("启动错误", zap.Error(err))
	}
}
