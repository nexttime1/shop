package handler

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	"option_service/proto"
)

func (o OptionServer) GetFavList(ctx context.Context, request *proto.UserFavRequest) (*proto.UserFavListResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (o OptionServer) AddUserFav(ctx context.Context, request *proto.UserFavRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (o OptionServer) DeleteUserFav(ctx context.Context, request *proto.UserFavRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (o OptionServer) GetUserFavDetail(ctx context.Context, request *proto.UserFavRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}
