package main

import (
	"encoding/json"
	"fmt"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dypnsapi20170525 "github.com/alibabacloud-go/dypnsapi-20170525/v3/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"strings"
)

// 全局变量：直接定义阿里云 Access Key
var (
	accessKeyID     = "LTAI5t8p6D7DSBazeKJq1inF"
	accessKeySecret = "oe1WX1CN7zxM6F7BvmPNpaTsHoxlY2"
)

func CreateClient() (_result *dypnsapi20170525.Client, _err error) {
	// ✅ 关键：直接设置 AccessKeyId 和 AccessKeySecret
	config := &openapi.Config{
		AccessKeyId:     tea.String(accessKeyID),     // 直接传字符串
		AccessKeySecret: tea.String(accessKeySecret), // 直接传字符串
		Endpoint:        tea.String("dypnsapi.aliyuncs.com"),
	}

	// 不需要 Credential 字段！
	_result, _err = dypnsapi20170525.NewClient(config)
	return _result, _err
}

func _main() (_err error) {
	client, _err := CreateClient()
	if _err != nil {
		return _err
	}

	sendSmsVerifyCodeRequest := &dypnsapi20170525.SendSmsVerifyCodeRequest{
		SignName:      tea.String("速通互联验证码"),
		TemplateCode:  tea.String("100001"),
		PhoneNumber:   tea.String("17564369039"),
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
		fmt.Printf("短信发送成功，响应结果：%s\n", util.ToJSONString(resp))
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

func main() {
	err := _main()
	if err != nil {
		panic(err)
	}
}
