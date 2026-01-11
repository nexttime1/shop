package api

import (
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"go.uber.org/zap"
	"user_web/common/res"
	"user_web/utils/aliyun"
)

type CaptchaApi struct{}

type CaptchaResponse struct {
	CaptchaId     string `json:"captcha_id"`
	CaptchaBase64 string `json:"captcha_base64"`
}

type SendSmsRequest struct {
	Mobile string `json:"mobile" binding:"required,mobile"`
	Role   string `json:"role" binding:"required,oneof=1 2 3"` // 1 代表 管理员 2代表普通用户 3代表游客
}

type VerifySmsRequest struct {
	Mobile string `json:"mobile" binding:"required,mobile"`
	Code   string `json:"code" binding:"required"`
}

var Store = base64Captcha.DefaultMemStore

func (api CaptchaApi) CaptchaCreateView(c *gin.Context) {
	// 生成 数字验证码 实例
	driver := base64Captcha.NewDriverDigit(80, 240, 5, 0.7, 80)
	//  绑定到
	captcha := base64Captcha.NewCaptcha(driver, Store)
	// 生成
	captchaId, base64s, answer, err := captcha.Generate()
	if err != nil {
		res.FailWithMsg(c, res.FailServiceCode, "生成验证码错误")
		return
	}
	zap.S().Infof("Answer : %s", answer)
	res.OkWithData(c, CaptchaResponse{
		CaptchaId:     captchaId,
		CaptchaBase64: base64s,
	})

}

func (api CaptchaApi) SendRegisterView(c *gin.Context) {
	var cr SendSmsRequest
	if err := c.ShouldBindJSON(&cr); err != nil {
		zap.S().Error(err)
		res.FailWithErr(c, res.FailArgumentCode, err)
		return
	}
	err := aliyun.SendCode(cr.Mobile)
	if err != nil {
		zap.S().Error(err)
		res.FailWithErr(c, res.FailServiceCode, err)
		return
	}
	res.OkWithMessage(c, "发送成功")

}

func (api CaptchaApi) VerifyCaptchaView(c *gin.Context) {
	var cr VerifySmsRequest
	if err := c.ShouldBindJSON(&cr); err != nil {
		res.FailWithErr(c, res.FailArgumentCode, err)
		return
	}
	success, err := aliyun.CheckSmsVerifyCode(cr.Mobile, cr.Code)
	if err != nil {
		zap.S().Errorf("错误： %v", err)
	}
	if success {
		res.OkWithMessage(c, "验证成功")
	} else {
		res.FailWithMsg(c, res.FailArgumentCode, "验证码错误")
	}

}
