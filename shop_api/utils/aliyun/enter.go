package aliyun

import (
	"encoding/json"
	"fmt"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dypnsapi20170525 "github.com/alibabacloud-go/dypnsapi-20170525/v3/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"shop_api/global"
	"strings"
)

func CreateClient() (_result *dypnsapi20170525.Client, _err error) {
	// 设置 AccessKeyId 和 AccessKeySecret
	config := &openapi.Config{
		AccessKeyId:     tea.String(global.Config.ALI.AccessKeyID),     // AccessKeyID
		AccessKeySecret: tea.String(global.Config.ALI.AccessKeySecret), // AccessKeySecret
		Endpoint:        tea.String("dypnsapi.aliyuncs.com"),
	}

	_result, _err = dypnsapi20170525.NewClient(config)
	return _result, _err
}

func SendCode(phone string) (_err error) {
	client, _err := CreateClient()
	if _err != nil {
		return _err
	}

	sendSmsVerifyCodeRequest := &dypnsapi20170525.SendSmsVerifyCodeRequest{
		SignName:      tea.String("速通互联验证码"),
		TemplateCode:  tea.String("100001"),
		PhoneNumber:   tea.String(phone),
		TemplateParam: tea.String("{\"code\":\"##code##\",\"min\":\"5\"}"),
	}
	runtime := &util.RuntimeOptions{}

	tryErr := func() (_e error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				_e = r
			}
		}()
		resp, _err := client.SendSmsVerifyCodeWithOptions(sendSmsVerifyCodeRequest, runtime)
		if _err != nil {
			return _err
		}
		fmt.Printf("短信发送成功，响应结果：%v\n", util.ToJSONString(resp))
		return nil
	}()

	if tryErr != nil {
		var error = &tea.SDKError{}
		if _t, ok := tryErr.(*tea.SDKError); ok {
			error = _t
		} else {
			error.Message = tea.String(tryErr.Error())
		}
		fmt.Println("短信发送失败，错误信息：", tea.StringValue(error.Message))
		var data interface{}
		d := json.NewDecoder(strings.NewReader(tea.StringValue(error.Data)))
		d.Decode(&data)
		if m, ok := data.(map[string]interface{}); ok {
			recommend, _ := m["Recommend"]
			fmt.Println("阿里云建议：", recommend)
		}
		_, _err = util.AssertAsString(error.Message)
		if _err != nil {
			return _err
		}
	}
	return _err
}
