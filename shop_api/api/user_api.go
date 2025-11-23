package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"shop_api/common/enum"
	"shop_api/common/res"
	"shop_api/connect"
	"shop_api/proto"
	"shop_api/service/user_service"
	"shop_api/utils/jwts"
)

type UserApi struct {
}

func (UserApi) UserListView(c *gin.Context) {
	client, conn, err := connect.UserConnectService(c)
	if err != nil {
		return
	}
	defer conn.Close()
	userListResponse, err := client.GetUserList(context.Background(), &proto.PageInfo{
		Page:  1,
		Limit: 5,
	})
	if err != nil {
		res.FailWithServiceMsg(c, err)
		return
	}
	res.OkWithList(c, userListResponse.Data, userListResponse.Total)

}

func (UserApi) UserLoginView(c *gin.Context) {

	var userLoginRequest user_service.UserLoginRequest
	err := c.ShouldBindJSON(&userLoginRequest)
	if err != nil {
		zap.S().Error(err)
		res.FailWithErr(c, res.FailArgumentCode, err)
		return
	}
	//验证码
	verifyResult := Store.Verify(userLoginRequest.CaptchaId, userLoginRequest.Answer, true)
	if !verifyResult {
		res.FailWithMsg(c, res.FailArgumentCode, "验证码错误")
		return
	}

	client, conn, err := connect.UserConnectService(c)
	if err != nil {
		return
	}
	defer conn.Close()
	// get User Info By Module
	ctx := context.Background()
	userInfo, err := client.GetUserMobile(ctx, &proto.MobileRequest{
		Mobile: userLoginRequest.Mobile,
	})
	if err != nil {
		res.FailWithServiceMsg(c, err)
		return
	}
	passwordResponse, err := client.CheckPassword(ctx, &proto.CheckPasswordReq{
		Password:          userLoginRequest.Password,
		EncryptedPassword: userInfo.Password,
	})
	if err != nil {
		res.FailWithServiceMsg(c, err)
		return
	}
	if passwordResponse.IsValid == false {
		res.FailWithMsg(c, res.FailArgumentCode, "密码错误")
		return
	}
	claim := jwts.Claims{
		UserID:   userInfo.Id,
		Username: userInfo.NickName,
		Role:     enum.RoleType(userInfo.Role),
	}
	token, err := jwts.GetToken(claim)
	if err != nil {
		res.FailWithServiceMsg(c, err)
		return
	}
	res.OkWithData(c, token)

}
