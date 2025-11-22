package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"shop_api/common/res"
	"shop_api/connect"
	"shop_api/proto"
	"shop_api/service"
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
	var userLoginRequest service.UserLoginRequest
	err := c.ShouldBindJSON(&userLoginRequest)
	if err != nil {
		zap.S().Error(err)
		res.FailWithErr(c, res.FailArgumentCode, err)
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
	res.OkWithMessage(c, "登录成功")

}
