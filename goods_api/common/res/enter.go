package res

import (
	"github.com/gin-gonic/gin"
	"goods_api/global"
	"goods_api/utils/validate"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type BaseResponse struct {
	Code Code   `json:"code"`
	Data any    `json:"data"`
	Msg  string `json:"msg"`
}

type DataListResponse struct {
	List  any   `json:"list"`
	Count int32 `json:"count"`
}
type Code int

const (
	SuccessCode      Code = 0    //成功
	FailLoginCode    Code = 1001 //登录错误
	FailServiceCode  Code = 1002 //服务异常
	FailArgumentCode Code = 1003
	NotFoundCode     Code = 1004
)

var empty = map[string]interface{}{}

func RPCErrorToHttp(err error) (Code, string) {
	if e, ok := status.FromError(err); ok {
		switch e.Code() {
		case codes.NotFound:
			return NotFoundCode, e.Message()
		case codes.InvalidArgument:
			return FailArgumentCode, e.Message()
		case codes.AlreadyExists:
			return FailArgumentCode, e.Message()
		case codes.Internal:
			return FailServiceCode, e.Message()
		case codes.Unavailable:
			return FailServiceCode, "服务不可用"
		default:
			return FailServiceCode, e.Message()
		}
	}
	return FailServiceCode, err.Error()
}

func Response(c *gin.Context, code Code, data interface{}, msg string) {
	c.JSON(200, BaseResponse{
		Code: code,
		Data: data,
		Msg:  msg,
	})
}

func OK(c *gin.Context, code Code, data interface{}, msg string) {
	Response(c, code, data, msg)
}

func OkWithMessage(c *gin.Context, msg string) {
	Response(c, SuccessCode, empty, msg)
}

func OkWithData(c *gin.Context, data interface{}) {
	Response(c, SuccessCode, data, "成功")
}

func OkWithList(c *gin.Context, list interface{}, count int32) {
	Response(c, SuccessCode, DataListResponse{
		List:  list,
		Count: count,
	}, "成功")

}

func Fail(c *gin.Context, code Code, data interface{}, msg string) {
	global.LevelFlag = true
	Response(c, code, data, msg)
}

func FailWithServiceMsg(c *gin.Context, err error) { //错误填在err里面
	code, msg := RPCErrorToHttp(err)
	global.LevelFlag = true
	Response(c, code, empty, msg)
}

func FailWithMsg(c *gin.Context, code Code, msg string) {
	global.LevelFlag = true
	Response(c, code, empty, msg)
}

func FailWithData(c *gin.Context, data interface{}) {
	global.LevelFlag = true
	Response(c, FailServiceCode, data, "失败")
}
func FailWithErr(c *gin.Context, code Code, err error) {
	data, _ := validate.ValidateErr(err)
	Fail(c, code, data, "失败")
}
