package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"order_web/global"
	"order_web/middleware"
	//myValidator "order_web/validator"
)

func Router() {
	gin.SetMode(global.Config.System.GinMode)
	r := gin.Default()

	// Gin 框架默认使用 go-playground/validator 这个库来做参数验证。  Engine : 返回底层实际的验证器引擎
	//validate, ok := binding.Validator.Engine().(*validator.Validate)
	//if ok { // Engine() 返回的是一个 interface{} 类型，我们需要把它转换成具体的 *validator.Validate 类型，才能调用它的方法
	//
	//	//规定签名必须是这个  func(fl validator.FieldLevel) bool
	//	_ = validate.RegisterValidation("mobile", myValidator.ValidateMobile)
	//}

	// 解决跨域问题
	r.Use(middleware.Cors())
	HealthRouter(r)
	ApiGroup := r.Group("/o/v1")
	OrderRouter(ApiGroup)
	CartRouter(ApiGroup)
	go func() {
		err := r.Run(global.Config.System.GetAddr())
		if err != nil {
			zap.L().Error("启动错误", zap.Error(err))
			return
		}
	}()

	zap.L().Info("router is running ...")
	// 上线使用这个
	//addr ,err := free_port.GetFreePort()
	//if err != nil {
	//	zap.L().Error(" Port 错误", zap.Error(err))
	//	return
	//}
	//err := r.Run(fmt.Sprintf("127.0.0.1:%d", addr))
	//if err != nil {
	//	zap.L().Error("启动错误", zap.Error(err))
	//}

}
