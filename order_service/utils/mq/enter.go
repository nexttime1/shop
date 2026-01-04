package mq

import (
	"context"
	"encoding/json"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"order_service/connect"
	"order_service/global"
	"order_service/models"
	"order_service/proto"
	"order_service/service"
)

type TransactionProducer struct {
	ID       int32
	Code     codes.Code
	Detail   string
	PriceSum float32
}

// When send transactional prepare(half) message succeed, this method will be invoked to execute local transaction.
func (t *TransactionProducer) ExecuteLocalTransaction(msg *primitive.Message) primitive.LocalTransactionState {
	var code codes.Code
	var detail string
	//var priceSum float32
	var request service.OrderTransitionRequest
	_ = json.Unmarshal(msg.Body, &request)
	//先拿到 选中的 good ID
	check := true
	var goodsId []int32
	var shopModels []models.ShoppingCartModel
	global.DB.Where(models.ShoppingCartModel{
		User:    request.UserId,
		Checked: &check,
	}).Find(&shopModels)
	if len(shopModels) == 0 {
		code = codes.NotFound
		detail = "请选择商品"
		return primitive.RollbackMessageState
	}
	goodNumMap := make(map[int32]int32)
	for _, shopModel := range shopModels {
		goodsId = append(goodsId, shopModel.Goods)
		goodNumMap[shopModel.Goods] = shopModel.Nums
	}

	// 调用good 微服务
	goodClient, conn, err := connect.GoodConnectService()
	if err != nil {
		zap.S().Error(err)
		t.Code = codes.Internal
		t.Detail = "服务启动失败"
		return primitive.RollbackMessageState
	}
	defer conn.Close()
	goods, err := goodClient.BatchGetGoods(context.Background(), &proto.BatchGoodsIdInfo{
		Id: goodsId,
	})
	if err != nil {
		zap.S().Error(err)
		t.Code = codes.Internal
		t.Detail = "商品查询失败"
		return primitive.RollbackMessageState
	}
	var PriceSum float32
	var orderGoods []*models.OrderGoodsModel
	var goodsInfo []*proto.GoodsInvInfo
	for _, goodModel := range goods.Data {
		PriceSum += goodModel.ShopPrice * float32(goodNumMap[goodModel.Id])
		orderGoods = append(orderGoods, &models.OrderGoodsModel{
			Goods:      goodModel.Id,
			GoodsName:  goodModel.Name,
			GoodsPrice: goodModel.ShopPrice,
			GoodImages: goodModel.GoodsFrontImage,
			Nums:       goodNumMap[goodModel.Id],
		})
		// 库存服务接收参数
		goodsInfo = append(goodsInfo, &proto.GoodsInvInfo{
			GoodsId: goodModel.Id,
			Num:     goodNumMap[goodModel.Id],
		})
	}
	// 预扣减库存
	inventoryClient, inventoryConn, err := connect.InventoryConnectService()
	if err != nil {
		zap.S().Error(err)
		t.Code = codes.Internal
		t.Detail = "服务启动失败"
		return primitive.RollbackMessageState
	}
	defer inventoryConn.Close()
	_, err = inventoryClient.Sell(context.Background(), &proto.SellInfo{GoodsInfo: goodsInfo, OrderSn: request.OrderSns})
	if err != nil {
		zap.S().Error(err)
		t.Code = codes.Internal
		t.Detail = "库存不足"
		return primitive.RollbackMessageState
	}

	// 生成订单表
	order := models.OrderModel{
		User:         request.UserId,
		OrderSn:      request.OrderSns,
		OrderMount:   PriceSum,
		Address:      request.Address,
		SignerName:   request.Name,
		SignerMobile: request.Mobile,
		Post:         request.Post,
	}
	//return primitive.CommitMessageState
	// 开启事务，保证操作原子性
	tx := global.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	err = tx.Create(&order).Error
	if err != nil {
		zap.S().Error(err)
		tx.Rollback()
		t.Code = codes.Internal
		t.Detail = "创建失败"
		return primitive.CommitMessageState
	}
	// 加上 订单ID
	for _, orderGood := range orderGoods {
		orderGood.Order = order.ID
	}
	// 生成 OrderGoodsModel 表数据
	err = tx.CreateInBatches(&orderGoods, 100).Error
	if err != nil {
		zap.S().Error(err)
		tx.Rollback()
		t.Code = codes.Internal
		t.Detail = "创建失败"
		return primitive.CommitMessageState
	}
	// 删除购物车中 已经生成订单的商品
	err = tx.Model(&models.ShoppingCartModel{}). // Model传空指针，指定操作shoppingcart表
		Where("user = ? AND checked = ?", request.UserId, check). // Where传查询条件
		Delete(&models.ShoppingCartModel{}).Error // Delete传指针（必须）
	if err != nil {
		zap.S().Error(err)
		tx.Rollback()
		t.Code = codes.Internal
		t.Detail = "删除失败"
		return primitive.CommitMessageState
	}
	t.Code = codes.OK
	t.ID = order.ID
	t.PriceSum = PriceSum
	// 发送延时消息  确保归还库存  发送普通消息就行
	p, err := rocketmq.NewProducer(producer.WithNameServer([]string{"192.168.163.132:9876"}))
	if err != nil {
		tx.Rollback()
		t.Code = codes.Internal
		t.Detail = "发送延时消息失败"
		return primitive.CommitMessageState
	}
	err = p.Start()
	if err != nil {
		tx.Rollback()
		t.Code = codes.Internal
		t.Detail = "发送延时消息失败"
		return primitive.CommitMessageState
	}
	message := primitive.NewMessage("order_timeout", msg.Body)
	message.WithDelayTimeLevel(3)
	_, err = p.SendSync(context.Background(), message)

	if err != nil {
		tx.Rollback()
		t.Code = codes.Internal
		t.Detail = "发送延时消息失败"
		return primitive.CommitMessageState
	}

	tx.Commit()
	return primitive.RollbackMessageState
}

// When no response to prepare(half) message. broker will send check message to check the transaction status, and this
// method will be invoked to get local transaction status.
func (t *TransactionProducer) CheckLocalTransaction(msg *primitive.MessageExt) primitive.LocalTransactionState {
	var request service.OrderTransitionRequest
	_ = json.Unmarshal(msg.Body, &request)
	var model models.OrderModel
	err := global.DB.Where("order_sn = ?", request.OrderSns).Take(&model).Error
	if err != nil {
		return primitive.CommitMessageState
	}

	return primitive.RollbackMessageState
}
