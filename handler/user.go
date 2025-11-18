package handler

import (
	"context"
	"github.com/sirupsen/logrus"
	"shop_service/common"
	"shop_service/models"
	"shop_service/proto"
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

//GetUserInfo(context.Context, *IdRequest) (*UserInfoResponse, error)
//GetUserMobile(context.Context, *MobileRequest) (*UserInfoResponse, error)
//CreateUser(context.Context, *CreateUserReq) (*UserInfoResponse, error)
//UpdateUser(context.Context, *UpdateUserReq) (*Response, error)
//CheckPassword(context.Context, *CheckPasswordReq) (*CheckPasswordResponse, error)
