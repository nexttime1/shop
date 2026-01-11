package api

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
	"user_web/common/enum"
	"user_web/common/res"
	"user_web/connect"
	"user_web/proto"
	"user_web/service/user_service"
	"user_web/utils/aliyun"
	"user_web/utils/jwts"
)

type UserApi struct {
}

func (UserApi) UserListView(c *gin.Context) {
	client, conn, err := connect.UserConnectService(c)
	if err != nil {
		return
	}
	defer conn.Close()
	userListResponse, err := client.GetUserList(context.WithValue(context.Background(), "ginContext", c), &proto.PageInfo{
		Page:  1,
		Limit: 5,
	})
	if err != nil {
		res.FailWithServiceMsg(c, err)
		return
	}
	var response []user_service.UserListResponse
	for _, v := range userListResponse.Data {
		password := v.Password[0:1] + "*****"
		loc, _ := time.LoadLocation("Asia/Shanghai")
		// 解析秒级时间戳 -> 转东八区时间
		birthTime := time.Unix(int64(v.BirthDay), 0).In(loc)
		birthDayStr := birthTime.Format("2006-01-02")
		response = append(response, user_service.UserListResponse{
			Id:       v.Id,
			Password: password,
			Mobile:   v.Mobile,
			NickName: v.NickName,
			BirthDay: birthDayStr,
			Gender:   v.Gender,
			Role:     int(v.Role),
		})
	}
	res.OkWithList(c, response, userListResponse.Total)

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
	ctx := context.WithValue(context.Background(), "ginContext", c)
	userInfo, err := client.GetUserMobile(ctx, &proto.MobileRequest{
		Mobile: userLoginRequest.Mobile,
	})
	if err != nil {
		res.FailWithServiceMsg(c, err)
		return
	}
	passwordResponse, err := client.CheckPassword(context.WithValue(context.Background(), "ginContext", c), &proto.CheckPasswordReq{
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

func (UserApi) UserRegisterView(c *gin.Context) {
	var cr user_service.UserRegisterRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		zap.S().Error(err)
		res.FailWithErr(c, res.FailArgumentCode, err)
		return
	}
	//  核实 验证码
	success, err := aliyun.CheckSmsVerifyCode(cr.Mobile, cr.Code)
	if !success {
		zap.S().Errorf("验证码问题success  %v", err)
		res.FailWithMsg(c, res.FailArgumentCode, "验证码错误")
		return
	}
	// 验证码对了
	client, conn, err := connect.UserConnectService(c)
	if err != nil {
		return
	}
	defer conn.Close()

	// 查一下 有没有这个手机号  链路追踪
	ctx := context.WithValue(context.Background(), "ginContext", c)
	userInfo, err := client.GetUserMobile(ctx, &proto.MobileRequest{
		Mobile: cr.Mobile,
	})
	if err == nil {
		_, err = client.UpdateUser(ctx, &proto.UpdateUserReq{
			Id:       userInfo.Id,
			Password: userInfo.Password,
		})
		if err != nil {
			res.FailWithServiceMsg(c, err)
			return
		}
		res.OkWithMessage(c, "您已经注册，密码更新成功")
		return
	}
	_, err = client.CreateUser(ctx, &proto.CreateUserReq{
		Password: cr.Password,
		NickName: fmt.Sprintf("user_%s", cr.Mobile),
		Mobile:   cr.Mobile,
		Role:     cr.Role,
	})
	if err != nil {
		res.FailWithServiceMsg(c, err)
		return
	}
	res.OkWithMessage(c, "创建成功")

}

func (UserApi) UserUpdateView(c *gin.Context) {
	// 先认证 看看你是不是 admin
	_claims, exist := c.Get("claims")
	if !exist {
		return
	}
	claims := _claims.(*jwts.MyClaims)
	var cr user_service.UserUpdateRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		zap.S().Error(err)
		res.FailWithErr(c, res.FailArgumentCode, err)
		return
	}
	UserCilent, conn, err := connect.UserConnectService(c)
	if err != nil {
		return
	}
	defer conn.Close()
	if cr.Id != claims.UserID {
		if claims.Role != enum.AdminRole {
			res.FailWithMsg(c, res.FailArgumentCode, "权限不够")
			return
		}
	}

	_, err = UserCilent.UpdateUser(context.Background(), &proto.UpdateUserReq{
		Id:       cr.Id,
		Password: cr.Password,
		NickName: cr.NickName,
		BirthDay: cr.BirthDay,
		Gender:   cr.Gender,
		Role:     cr.Role,
		UserId:   claims.UserID,
	})
	if err != nil {
		zap.S().Error(err)
		res.FailWithServiceMsg(c, err)
		return
	}

	res.OkWithMessage(c, "更新成功")

}
