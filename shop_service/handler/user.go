package handler

import (
	"context"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"shop_service/common"
	"shop_service/global"
	"shop_service/models"
	"shop_service/proto"
	"time"
)

type UserSever struct {
}

func UserResponseFunction(user models.UserModel) proto.UserInfoResponse {
	response := proto.UserInfoResponse{
		Id:       user.ID,
		Mobile:   user.Mobile,
		Password: user.Password,
		NickName: user.NickName,
		Gender:   user.Gender,
		Role:     int32(user.Role),
	}
	if user.Birthday != nil {
		response.BirthDay = uint64(user.Birthday.Unix())
	}

	return response
}

func (UserSever) GetUserList(c context.Context, pageInfo *proto.PageInfo) (*proto.UserListResponse, error) {
	list, count, err := common.ListQuery(models.UserModel{}, common.Options{
		PageInfo: common.PageInfo{
			Limit: pageInfo.Limit,
			Page:  pageInfo.Page,
		},
	})
	if err != nil {
		logrus.Errorf("get user list error: %v", err)
		return nil, err
	}
	var userList []*proto.UserInfoResponse
	for _, user := range list {
		response := UserResponseFunction(user)
		userList = append(userList, &response)
	}

	return &proto.UserListResponse{
		Total: int32(count),
		Data:  userList,
	}, nil

}

func (UserSever) GetUserInfo(ctx context.Context, id *proto.IdRequest) (*proto.UserInfoResponse, error) {
	var user models.UserModel
	count := global.DB.Where("id = ?", id.Id).First(&user).RowsAffected
	if count != 1 {
		return nil, status.Error(codes.NotFound, "用户不存在")
	}
	response := UserResponseFunction(user)
	return &response, nil
}

func (UserSever) GetUserMobile(ctx context.Context, mobile *proto.MobileRequest) (*proto.UserInfoResponse, error) {
	var user models.UserModel
	result := global.DB.Where("mobile = ?", mobile).First(&user)
	if result.RowsAffected == 0 {
		return nil, status.Error(codes.NotFound, "用户不存在")
	}
	if result.Error != nil {
		return nil, result.Error
	}
	response := UserResponseFunction(user)
	return &response, nil

}
func (UserSever) CreateUser(ctx context.Context, req *proto.CreateUserReq) (*proto.UserInfoResponse, error) {
	var user models.UserModel
	count := global.DB.Where("mobile = ?", req.Mobile).First(&user).RowsAffected
	if count != 0 {
		return nil, status.Error(codes.AlreadyExists, "该手机号已经注册")
	}
	user.Mobile = req.Mobile
	user.Password = req.Password
	user.NickName = req.NickName
	err := global.DB.Create(&user).Error
	if err != nil {
		logrus.Errorf("create user error: %v", err)
		return nil, status.Error(codes.Internal, "创建用户失败")
	}
	response := UserResponseFunction(user)
	return &response, nil

}

func (UserSever) UpdateUser(ctx context.Context, req *proto.UpdateUserReq) (*proto.Response, error) {
	var user models.UserModel
	count := global.DB.Where("id = ?", req.Id).First(&user).RowsAffected
	if count == 0 {
		return nil, status.Error(codes.NotFound, "用户不存在")
	}
	BirthDay := time.Unix(int64(req.BirthDay), 0)
	user.Birthday = &BirthDay
	user.NickName = req.NickName
	user.Gender = req.Gender
	user.Role = int(req.Role)

	err := global.DB.Save(user).Error
	if err != nil {
		logrus.Errorf("update user error: %v", err)
		return nil, status.Error(codes.Internal, "用户更新失败")
	}
	response := proto.Response{
		Code: int32(codes.OK),
		Msg:  "更新成功",
	}
	return &response, nil

}
func (UserSever) CheckPassword(ctx context.Context, check *proto.CheckPasswordReq) (*proto.CheckPasswordResponse, error) {
	return nil, status.Error(codes.Unimplemented, "wait")
}
