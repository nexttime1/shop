package handler

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"option_service/common"
	"option_service/connect"
	"option_service/global"

	"option_service/models"
	"option_service/proto"
)

func (o OptionServer) GetFavList(ctx context.Context, request *proto.UserFavRequest) (*proto.UserFavListResponse, error) {
	parentSpan := opentracing.SpanFromContext(ctx)
	// 链路记录
	mysqlSpan := opentracing.GlobalTracer().StartSpan("mysql_option", opentracing.ChildOf(parentSpan.Context()))
	list, count, err := common.ListQuery(models.UserCollectionModel{UserId: request.UserId, GoodId: request.GoodsId}, common.Options{})
	if err != nil {
		zap.S().Error(err)
		return nil, status.Errorf(codes.Internal, "查询失败")
	}
	mysqlSpan.Finish()
	var FavModels []*proto.UserFavResponse
	for _, model := range list {
		FavModels = append(FavModels, &proto.UserFavResponse{
			UserId:  model.UserId,
			GoodsId: model.GoodId,
		})
	}
	response := &proto.UserFavListResponse{
		Total: int32(count),
		Data:  FavModels,
	}
	return response, nil

}

func (o OptionServer) AddUserFav(ctx context.Context, request *proto.UserFavRequest) (*emptypb.Empty, error) {
	parentSpan := opentracing.SpanFromContext(ctx)
	// 链路记录
	mysqlSpan := opentracing.GlobalTracer().StartSpan("mysql_option", opentracing.ChildOf(parentSpan.Context()))
	goodClient, conn, err := connect.GoodConnectService()
	if err != nil {
		zap.S().Error(err)
		return nil, status.Errorf(codes.Internal, "连接失败")
	}
	defer conn.Close()
	_, err = goodClient.GetGoodsDetail(context.Background(), &proto.GoodInfoRequest{
		Id: request.GoodsId,
	})
	if err != nil {
		zap.S().Error(err)
		return nil, status.Errorf(codes.Internal, "商品不存在")
	}

	model := models.UserCollectionModel{
		UserId: request.UserId,
		GoodId: request.GoodsId,
	}
	err = global.DB.Create(&model).Error
	if err != nil {
		zap.S().Error(err)
		return nil, status.Errorf(codes.Internal, "创建失败")
	}
	mysqlSpan.Finish()
	return &emptypb.Empty{}, nil

}

func (o OptionServer) DeleteUserFav(ctx context.Context, request *proto.UserFavRequest) (*emptypb.Empty, error) {
	parentSpan := opentracing.SpanFromContext(ctx)
	// 链路记录
	mysqlSpan := opentracing.GlobalTracer().StartSpan("mysql_option", opentracing.ChildOf(parentSpan.Context()))
	count := global.DB.Unscoped().Where("good_id = ? and user_id = ?", request.GoodsId, request.UserId).Delete(&models.UserCollectionModel{}).RowsAffected
	if count == 0 {
		return nil, status.Errorf(codes.NotFound, "未找到")
	}
	mysqlSpan.Finish()
	return &emptypb.Empty{}, nil
}

func (o OptionServer) GetUserFavDetail(ctx context.Context, request *proto.UserFavRequest) (*emptypb.Empty, error) {
	parentSpan := opentracing.SpanFromContext(ctx)
	// 链路记录
	mysqlSpan := opentracing.GlobalTracer().StartSpan("mysql_option", opentracing.ChildOf(parentSpan.Context()))
	var model models.UserCollectionModel
	count := global.DB.Debug().Where("good_id = ? and user_id = ?", request.GoodsId, request.UserId).Take(&model).RowsAffected
	if count == 0 {
		return nil, status.Errorf(codes.NotFound, "未找到")
	}
	mysqlSpan.Finish()
	return &emptypb.Empty{}, nil
}
