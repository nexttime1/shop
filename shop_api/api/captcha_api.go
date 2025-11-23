package api

import (
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"go.uber.org/zap"
	"shop_api/common/res"
)

type CaptchaApi struct{}

type CaptchaResponse struct {
	CaptchaId     string `json:"captcha_id"`
	CaptchaBase64 string `json:"captcha_base64"`
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
