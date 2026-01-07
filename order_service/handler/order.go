package handler

import (
	"context"
	"encoding/json"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/opentracing/opentracing-go"
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

var GlobalOrderServer *OrderSever

type OrderSever struct {
	transactionProducer rocketmq.TransactionProducer // 复用的生产者实例
	MessageProducer     rocketmq.Producer            //复用 普通消息的生产者实例
	// 事务监听器需要考虑线程安全，监听器也作为成员变量
	orderListener *mq.TransactionProducer // 监听器 便于复用
}

// InitProducer 初始化方法：在服务启动时调用 只执行一次
func (o *OrderSever) InitProducer() error {
	// 初始化事务监听器
	o.orderListener = &mq.TransactionProducer{}

	// 初始化延时消息生产者
	if err := o.orderListener.InitDelayProducer(); err != nil {
		zap.L().Error("初始化延时消息生产者失败", zap.Error(err))
		return err
	}
	// 创建事务生产者（只执行一次）
	producerIns, err := rocketmq.NewTransactionProducer(
		o.orderListener,
		producer.WithNameServer([]string{global.Config.RocketMQ.Addr()}),
		producer.WithGroupName(global.Config.RocketMQ.GroupName),
	)
	if err != nil {
		zap.L().Error("RocketMQ创建事务生产者失败", zap.Error(err))
		return err
	}

	// 启动生产者
	if err = producerIns.Start(); err != nil {
		zap.L().Error("启动生产者错误", zap.Error(err))
		return err
	}

	// 赋值给成员变量，供后续复用
	o.transactionProducer = producerIns

	// 初始化 普通消息
	p, err := rocketmq.NewProducer(producer.WithNameServer([]string{global.Config.RocketMQ.Addr()}))
	if err != nil {
		zap.L().Error("RocketMQ创建普通生产者失败", zap.Error(err))
		return err
	}
	err = p.Start()
	if err != nil {
		zap.L().Error("启动生产者错误", zap.Error(err))
		return err
	}
	o.MessageProducer = p
	return nil
}

// CloseProducer 程序退出时调用，释放资源
func (o *OrderSever) CloseProducer() error {
	// 关闭事务生产者
	if o.transactionProducer != nil {
		if err := o.transactionProducer.Shutdown(); err != nil {
			zap.L().Error("关闭事务生产者失败", zap.Error(err))
			return err
		}
	}
	// 关闭普通生产者
	if o.MessageProducer != nil {
		if err := o.MessageProducer.Shutdown(); err != nil {
			zap.L().Error("关闭普通生产者失败", zap.Error(err))
			return err
		}
	}

	// 关闭延时消息生产者
	if o.orderListener != nil {
		if err := o.orderListener.CloseDelayProducer(); err != nil {
			zap.L().Error("关闭延时消息生产者失败", zap.Error(err))
			return err
		}
	}
	return nil
}
func (o *OrderSever) CreateOrder(ctx context.Context, request *proto.OrderRequest) (*proto.OrderInfoResponse, error) {
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
	msg := primitive.NewMessage(global.Config.RocketMQ.Topic, data)
	//half 消息 如果回复了 我就调用本地事务 也就是 ExecuteLocalTransaction
	parentSpan := opentracing.SpanFromContext(ctx)
	if parentSpan != nil {
		// 创建 carrier 用于存储 span 上下文（用 message 属性作为载体）
		carrier := make(opentracing.TextMapCarrier)
		// 将 span 上下文注入到 carrier
		err := opentracing.GlobalTracer().Inject(
			parentSpan.Context(),
			opentracing.TextMap,
			carrier,
		)
		if err == nil {
			// 将 carrier 中的键值对存入 message 属性
			for k, v := range carrier {
				msg.WithProperty(k, v)
			}
		} else {
			zap.S().Warn("注入链路上下文到消息失败", zap.Error(err))
		}
	}
	// 链路记录
	halfSpan := opentracing.GlobalTracer().StartSpan("发送 half消息", opentracing.ChildOf(parentSpan.Context()))
	_, err := o.transactionProducer.SendMessageInTransaction(ctx, msg)
	if err != nil {
		zap.S().Error(err)
		return nil, err
	}
	halfSpan.Finish()

	//根据 TransactionId 获取请求状态
	statusInfo, ok := o.orderListener.GetTransactionStatus(model.OrderSns)
	if !ok {
		zap.S().Error("获取事务状态失败", zap.String("txId", model.OrderSns))
		return nil, status.Error(codes.Internal, "获取订单状态失败")
	}

	if statusInfo.Code != codes.OK {
		return nil, status.Error(statusInfo.Code, statusInfo.Detail)
	}

	return &proto.OrderInfoResponse{Id: statusInfo.ID, OrderSn: model.OrderSns, Total: statusInfo.PriceSum}, nil

}

func (o *OrderSever) OrderList(ctx context.Context, request *proto.OrderFilterRequest) (*proto.OrderListResponse, error) {
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

func (o *OrderSever) OrderDetail(ctx context.Context, request *proto.OrderRequest) (*proto.OrderInfoDetailResponse, error) {
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

func (o *OrderSever) UpdateOrderStatus(ctx context.Context, req *proto.OrderStatus) (*emptypb.Empty, error) {
	result := global.DB.Model(&models.OrderModel{}).Where("order_sn = ?", req.OrderSn).Update("status", req.Status)
	if result.Error != nil || result.RowsAffected == 0 {
		return nil, status.Errorf(codes.Internal, "订单不存在")
	}

	return &emptypb.Empty{}, nil
}
