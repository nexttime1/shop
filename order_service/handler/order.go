package handler

import (
	"context"
	"encoding/json"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"order_service/common"

	"order_service/global"
	"order_service/models"
	"order_service/proto"
	"order_service/service"
	"order_service/utils/mq"
)

type OrderSever struct {
}

func (o OrderSever) CreateOrder(ctx context.Context, request *proto.OrderRequest) (*proto.OrderInfoResponse, error) {
	model := service.OrderTransitionRequest{
		Id:       request.Id,
		UserId:   request.UserId,
		Address:  request.Address,
		Name:     request.Name,
		Mobile:   request.Mobile,
		Post:     request.Post,
		OrderSns: service.RandomSns(request.UserId),
	}
	data, _ := json.Marshal(model)
	var orderListener mq.TransactionProducer
	transactionProducer, err := rocketmq.NewTransactionProducer(
		&orderListener,
		producer.WithNameServer([]string{"192.168.163.132:9876"}),
	)
	if err != nil {
		panic(err)
	}
	err = transactionProducer.Start()
	if err != nil {
		panic(err)
	}
	_, err = transactionProducer.SendMessageInTransaction(context.Background(), primitive.NewMessage("shop_reback", data))
	if err != nil {
		zap.S().Error(err)
		return nil, err
	}
	if orderListener.Code != codes.OK {
		return nil, status.Error(orderListener.Code, orderListener.Detail)
	}

	return &proto.OrderInfoResponse{Id: orderListener.ID, OrderSn: model.OrderSns, Total: orderListener.PriceSum}, nil

}

func (o OrderSever) OrderList(ctx context.Context, request *proto.OrderFilterRequest) (*proto.OrderListResponse, error) {
	// 管理员看所有的列表   而用户看自己的  区别是 看web端给我发不发id
	response := &proto.OrderListResponse{}
	pageInfo := common.PageInfo{
		Page:  request.PageNum,
		Limit: request.PageSize,
	}
	var err error
	var list []models.OrderModel
	var count int
	if request.UserId == 0 {
		// 管理员看所有
		list, count, err = common.ListQuery(models.OrderModel{}, common.Options{
			PageInfo: pageInfo,
		})
		if err != nil {
			zap.S().Error(err)
			return nil, status.Error(codes.Internal, "查询错误")
		}

	} else {
		list, count, err = common.ListQuery(models.OrderModel{User: request.UserId}, common.Options{
			PageInfo: pageInfo,
		})
		if err != nil {
			zap.S().Error(err)
			return nil, status.Error(codes.Internal, "查询错误")
		}
	}
	response.Total = int32(count)
	var modelsInfo []*proto.OrderInfoResponse
	for _, item := range list {

		modelsInfo = append(modelsInfo, &proto.OrderInfoResponse{
			Id:      item.ID,
			UserId:  item.User,
			OrderSn: item.OrderSn,
			PayType: item.PayType,
			Status:  item.Status,
			Post:    item.Post,
			Total:   item.OrderMount,
			Address: item.Address,
			Name:    item.SignerName,
			Mobile:  item.SignerMobile,
		})
	}
	response.Data = modelsInfo
	return response, nil
}

func (o OrderSever) OrderDetail(ctx context.Context, request *proto.OrderRequest) (*proto.OrderInfoDetailResponse, error) {
	// 如果传userId  那就查这个用户的  不传就是全部的
	response := &proto.OrderInfoDetailResponse{}
	var model models.OrderModel
	err := global.DB.Where(models.OrderModel{User: request.UserId, Model: models.Model{ID: request.Id}}).Take(&model).Error
	if err != nil {
		zap.S().Error(err)
		return nil, status.Error(codes.NotFound, "订单不存在")
	}
	response.OrderInfo = &proto.OrderInfoResponse{
		Id:      model.ID,
		UserId:  model.User,
		OrderSn: model.OrderSn,
		PayType: model.PayType,
		Status:  model.Status,
		Post:    model.Post,
		Total:   model.OrderMount,
		Address: model.Address,
		Name:    model.SignerName,
		Mobile:  model.SignerMobile,
	}
	// 找一下商品
	var goodModels []models.OrderGoodsModel
	global.DB.Where("`order` = ?", model.ID).Find(&goodModels)
	var Goods []*proto.OrderItemResponse
	for _, item := range goodModels {
		Goods = append(Goods, &proto.OrderItemResponse{
			Id:         item.ID,
			OrderId:    item.Order,
			GoodsId:    item.Goods,
			GoodsName:  item.GoodsName,
			GoodsPrice: item.GoodsPrice,
			Nums:       item.Nums,
		})
	}
	response.Goods = Goods
	return response, nil

}

func (o OrderSever) UpdateOrderStatus(ctx context.Context, req *proto.OrderStatus) (*emptypb.Empty, error) {
	result := global.DB.Model(&models.OrderModel{}).Where("order_sn = ?", req.OrderSn).Update("status", req.Status)
	if result.Error != nil || result.RowsAffected == 0 {
		return nil, status.Errorf(codes.Internal, "订单不存在")
	}

	return &emptypb.Empty{}, nil
}
