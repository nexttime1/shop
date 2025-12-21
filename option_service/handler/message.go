package handler

import (
	"context"
	"option_service/proto"
)

func (o OptionServer) MessageList(ctx context.Context, request *proto.MessageRequest) (*proto.MessageListResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (o OptionServer) CreateMessage(ctx context.Context, request *proto.MessageRequest) (*proto.MessageResponse, error) {
	//TODO implement me
	panic("implement me")
}
