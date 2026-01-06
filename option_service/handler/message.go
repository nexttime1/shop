package handler

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"option_service/common"
	"option_service/global"
	"option_service/models"
	"option_service/models/enum"
	"option_service/proto"
)

func (o OptionServer) MessageList(ctx context.Context, request *proto.MessageRequest) (*proto.MessageListResponse, error) {
	parentSpan := opentracing.SpanFromContext(ctx)
	// 链路记录
	mysqlSpan := opentracing.GlobalTracer().StartSpan("mysql_option", opentracing.ChildOf(parentSpan.Context()))
	var MessageModels []*proto.MessageResponse
	list, count, err := common.ListQuery(models.LeavingMessageModel{UserId: request.UserId}, common.Options{})
	if err != nil {
		zap.S().Error(err)
		return nil, status.Error(codes.Internal, "查询错误")
	}
	for _, model := range list {
		MessageModels = append(MessageModels, &proto.MessageResponse{
			Id:          model.ID,
			UserId:      model.UserId,
			MessageType: int32(model.MessageType),
			Subject:     model.Subject,
			Message:     model.Message,
			File:        model.File,
		})
	}
	mysqlSpan.Finish()
	response := &proto.MessageListResponse{
		Total: int32(count),
		Data:  MessageModels,
	}
	return response, nil

}

func (o OptionServer) CreateMessage(ctx context.Context, request *proto.MessageRequest) (*proto.MessageResponse, error) {
	parentSpan := opentracing.SpanFromContext(ctx)
	// 链路记录
	mysqlSpan := opentracing.GlobalTracer().StartSpan("mysql_option", opentracing.ChildOf(parentSpan.Context()))
	model := models.LeavingMessageModel{
		UserId:      request.UserId,
		MessageType: enum.MessageType(request.MessageType),
		Subject:     request.Subject,
		Message:     request.Message,
		File:        request.File,
	}
	err := global.DB.Create(&model).Error

	if err != nil {
		zap.S().Error(err)
		return nil, status.Error(codes.Internal, "创建失败")
	}
	mysqlSpan.Finish()
	response := &proto.MessageResponse{
		Id: model.ID,
	}
	return response, nil
}
