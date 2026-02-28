package mq

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"order_service/connect"
	"order_service/global"
	"order_service/models"
	"order_service/proto"
	"order_service/service"
	"sync"
	"time"
)

// TransactionStatus 定义每个事务请求的状态（非单例）
type TransactionStatus struct {
	ID       int32
	Code     codes.Code
	Detail   string
	PriceSum float32
}

// TransactionProducer 事务监听器：用并发安全Map存储请求状态
type TransactionProducer struct {
	// 并发安全Map：key=消息TransactionId，value=该请求的状态
	statusMap sync.Map
	// 复用的延时消息生产者
	delayProducer rocketmq.Producer
}

// InitDelayProducer 初始化延时消息生产者（只执行一次）
func (t *TransactionProducer) InitDelayProducer() error {
	if t.delayProducer != nil {
		return nil
	}
	p, err := rocketmq.NewProducer(
		producer.WithNameServer([]string{global.Config.RocketMQ.Addr()}),
		// 延时消息 生产者组
		producer.WithGroupName(global.Config.RocketMQ.DelayGroupName),
	)
	if err != nil {
		return err
	}
	if err = p.Start(); err != nil {
		return err
	}
	t.delayProducer = p
	return nil
}

// CloseDelayProducer 关闭延时消息生产者
func (t *TransactionProducer) CloseDelayProducer() error {
	if t.delayProducer != nil {
		return t.delayProducer.Shutdown()
	}
	return nil
}

// When send transactional prepare(half) message succeed, this method will be invoked to execute local transaction.
func (t *TransactionProducer) ExecuteLocalTransaction(msg *primitive.Message) primitive.LocalTransactionState {
	// 初始化当前请求的状态
	status := &TransactionStatus{
		Code:   codes.Internal,
		Detail: "未知错误",
	}
	// 链路追踪
	var span opentracing.Span
	// 从 message 属性创建 carrier
	carrier := make(opentracing.TextMapCarrier)
	for k, v := range msg.GetProperties() {
		carrier.Set(k, v)
	}
	// 提取 span 上下文
	parentCtx, err := opentracing.GlobalTracer().Extract(
		opentracing.TextMap,
		carrier,
	)
	if err != nil {
		// 提取失败则创建根 span
		span = opentracing.GlobalTracer().StartSpan("ExecuteLocalTransaction")
	} else {
		// 提取成功则创建子 span
		span = opentracing.GlobalTracer().StartSpan("ExecuteLocalTransaction", opentracing.ChildOf(parentCtx))
	}

	ctx := opentracing.ContextWithSpan(context.Background(), span) // 为调用微服务传递 跟踪链路

	local_prepare := opentracing.GlobalTracer().StartSpan("local_prepare", opentracing.ChildOf(span.Context()))
	// 关联消息ID和状态（TransactionId是每个事务消息的唯一标识）
	var request service.OrderTransitionRequest
	_ = json.Unmarshal(msg.Body, &request)
	transactionId := request.OrderSns
	zap.S().Infof("transactionId为: %s", transactionId)
	// 空值保护
	if transactionId == "" {
		zap.S().Warn("事务消息 TransactionId 为空")
		status.Code = codes.Internal
		status.Detail = "事务标识为空"
		return primitive.RollbackMessageState
	}
	history := models.OrderStockHistory{
		OrderSn: transactionId,
		Status:  0,
	}
	global.DB.Create(&history) // 不会错
	t.statusMap.Store(transactionId, status)
	//var priceSum float32

	//先拿到 选中的 good ID
	check := true
	var goodsId []int32
	var shopModels []models.ShoppingCartModel
	global.DB.Where(models.ShoppingCartModel{
		User:    request.UserId,
		Checked: &check,
	}).Find(&shopModels)
	if len(shopModels) == 0 {
		status.Code = codes.NotFound
		status.Detail = "请选择商品"
		return primitive.RollbackMessageState
	}
	goodNumMap := make(map[int32]int32)
	for _, shopModel := range shopModels {
		goodsId = append(goodsId, shopModel.Goods)
		goodNumMap[shopModel.Goods] = shopModel.Nums
	}
	local_prepare.Finish()
	// 调用good 微服务
	goodService := opentracing.GlobalTracer().StartSpan("good_service", opentracing.ChildOf(span.Context()))
	goodClient, conn, err := connect.GoodConnectService()
	if err != nil {
		zap.S().Error(err)
		status.Code = codes.Internal
		status.Detail = "服务启动失败"
		return primitive.RollbackMessageState
	}
	defer conn.Close()
	goods, err := goodClient.BatchGetGoods(ctx, &proto.BatchGoodsIdInfo{
		Id: goodsId,
	})
	if err != nil {
		zap.S().Error(err)
		status.Code = codes.Internal
		status.Detail = "商品查询失败"
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
	goodService.Finish()

	// 预扣减库存
	stockService := opentracing.GlobalTracer().StartSpan("stock_service", opentracing.ChildOf(span.Context()))
	inventoryClient, inventoryConn, err := connect.InventoryConnectService()
	if err != nil {
		zap.S().Error(err)
		status.Code = codes.Internal
		status.Detail = "服务启动失败"
		return primitive.RollbackMessageState
	}
	defer inventoryConn.Close()
	_, err = inventoryClient.Sell(ctx, &proto.SellInfo{GoodsInfo: goodsInfo, OrderSn: request.OrderSns})
	if err != nil {
		zap.S().Error(err)
		status.Code = codes.Internal
		status.Detail = "库存不足"
		return primitive.RollbackMessageState
	}
	stockService.Finish()
	// 这个时候 去修改一下 history
	localMysql := opentracing.GlobalTracer().StartSpan("update_local_mysql", opentracing.ChildOf(span.Context()))
	history.Status = 1
	global.DB.Save(&history)
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
		status.Code = codes.Internal
		status.Detail = "创建失败"
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
		status.Code = codes.Internal
		status.Detail = "创建失败"
		return primitive.CommitMessageState
	}
	localMysql.Finish()
	DeleteCart := opentracing.GlobalTracer().StartSpan("delete_shop_cart", opentracing.ChildOf(span.Context()))
	// 删除购物车中 已经生成订单的商品
	err = tx.Model(&models.ShoppingCartModel{}). // Model传空指针，指定操作shoppingcart表
							Where("user = ? AND checked = ?", request.UserId, check). // Where传查询条件
							Delete(&models.ShoppingCartModel{}).Error                 // Delete传指针（必须）
	if err != nil {
		zap.S().Error(err)
		tx.Rollback()
		status.Code = codes.Internal
		status.Detail = "删除失败"
		return primitive.CommitMessageState
	}
	DeleteCart.Finish()
	// 发送延时消息  确保归还库存  发送普通消息就行   复用生产者
	delayMessage := opentracing.GlobalTracer().StartSpan("send_delay_message", opentracing.ChildOf(span.Context()))
	delayMsg := primitive.NewMessage(global.Config.RocketMQ.DelayTopic, msg.Body)
	delayMsg.WithDelayTimeLevel(6) // 延时级别6（根据RocketMQ配置对应时间）
	if _, err = t.delayProducer.SendSync(context.Background(), delayMsg); err != nil {
		zap.S().Error("发送延时消息失败", zap.Error(err))
		tx.Rollback()
		status.Detail = "发送延时消息失败"
		return primitive.CommitMessageState
	}
	delayMessage.Finish()

	err = tx.Commit().Error
	if err != nil {
		zap.S().Error("事务提交错误", zap.Error(err))
		status.Code = codes.Internal
		status.Detail = "事务提交错误"
		return primitive.CommitMessageState
	}

	status.Code = codes.OK
	status.ID = order.ID
	status.PriceSum = PriceSum
	return primitive.RollbackMessageState
}

// CheckLocalTransaction When no response to prepare(half) message. broker will send check message to check the transaction status, and this
// method will be invoked to get local transaction status.
func (t *TransactionProducer) CheckLocalTransaction(msg *primitive.MessageExt) primitive.LocalTransactionState {
	// 链路追踪
	var span opentracing.Span
	// 从 message 属性创建 carrier
	carrier := make(opentracing.TextMapCarrier)
	for k, v := range msg.GetProperties() {
		carrier.Set(k, v)
	}
	// 提取 span 上下文
	parentCtx, err := opentracing.GlobalTracer().Extract(
		opentracing.TextMap,
		carrier,
	)
	if err != nil {
		// 提取失败则创建根 span
		span = opentracing.GlobalTracer().StartSpan("CheckLocalTransaction")
	} else {
		// 提取成功则创建子 span
		span = opentracing.GlobalTracer().StartSpan("CheckLocalTransaction", opentracing.ChildOf(parentCtx))
	}
	defer span.Finish()
	// 先拿出 事务ID 如果 都执行成功 本地事务 扣减库存 全部完成 但是MQ 没收到最后的 return 首先判断 Code == OK
	id := msg.TransactionId
	status, ok := t.GetTransactionStatus(id)
	if !ok {
		zap.S().Errorf("id 为: %v", status)
		zap.S().Errorf("怎么可能走到这里")
		return primitive.RollbackMessageState
	}
	if status.Code == codes.OK {
		return primitive.RollbackMessageState
	}
	// 这里说明遇到问题了  需要看看 到底是 扣没扣减库存
	var history models.OrderStockHistory
	err = global.DB.Where("order_sn = ?", id).Take(&history).Error
	if err != nil {
		zap.S().Errorf("不可能错 %v", err)
		return primitive.RollbackMessageState
	}
	if history.Status == 0 {
		return primitive.RollbackMessageState
	}

	var request service.OrderTransitionRequest
	_ = json.Unmarshal(msg.Body, &request)
	var model models.OrderModel
	err = global.DB.Where("order_sn = ?", request.OrderSns).Take(&model).Error
	if err != nil {
		return primitive.CommitMessageState
	}

	return primitive.RollbackMessageState
}

// GetTransactionStatus 根据TransactionId获取请求状态（供外层调用）
func (t *TransactionProducer) GetTransactionStatus(transactionId string) (*TransactionStatus, bool) {
	val, ok := t.statusMap.Load(transactionId)
	if !ok {
		return nil, false
	}
	status, ok := val.(*TransactionStatus)
	return status, ok
}

// DeleteTransactionStatus 清理已完成的事务状态（避免内存泄漏）
func (t *TransactionProducer) DeleteTransactionStatus(transactionId string) {
	t.statusMap.Delete(transactionId)
}

// 发送延迟消息  带重试
func (t *TransactionProducer) sendDelayMsgWithRetry(msg *primitive.Message, orderSn string) (*primitive.SendResult, error) {
	var (
		sendResult *primitive.SendResult
		sendErr    error
		delayMsg   = primitive.NewMessage(global.Config.RocketMQ.DelayTopic, msg.Body)
	)
	delayMsg.WithDelayTimeLevel(16) // 30分钟延时

	var retryIdx int32
	// 执行重试逻辑
	for retryIdx = 0; retryIdx < global.Config.RocketMQ.MaxRetryTimes; retryIdx++ {
		// 发送消息
		sendResult, sendErr = t.delayProducer.SendSync(context.Background(), delayMsg)

		// 构建通用日志字段
		logFields := []zap.Field{
			zap.String("order_sn", orderSn),
			zap.Int32("retry_times", retryIdx+1), // 重试次数（1=首次，2=第1次重试，3=第2次重试）
			zap.Bool("send_success", sendErr == nil),
		}
		if sendResult != nil {
			logFields = append(logFields,
				zap.String("msg_id", sendResult.MsgID),
				zap.String("send_status", SendStatusToString(sendResult.Status)),
			)
		}
		if sendErr != nil {
			logFields = append(logFields, zap.Error(sendErr))
		}

		// 发送成功：记录日志并返回
		if sendResult != nil && sendErr == nil && sendResult.Status == primitive.SendOK {
			zap.S().Info("延时消息发送成功", logFields)
			return sendResult, nil
		}

		// 发送失败：记录日志，判断是否继续重试
		if retryIdx == global.Config.RocketMQ.MaxRetryTimes-1 {
			// 最后一次重试失败：记录错误日志（标记最终失败）
			zap.S().Error("延时消息3次发送均失败", logFields)
		} else {
			// 非最后一次：记录警告日志，等待后重试
			zap.S().Warn("延时消息发送失败，准备重试", logFields)
			// 递增间隔重试（避免高频重试）：第1次等500ms，第2次等1000ms
			baseRetryDelay := time.Duration(global.Config.RocketMQ.BaseRetryDelay) * time.Millisecond

			// 2. 再计算递增的重试间隔
			retryDelay := baseRetryDelay * time.Duration(retryIdx+1)
			time.Sleep(retryDelay)
		}
	}

	// 3次重试均失败：返回最终结果
	return sendResult, sendErr
}

// SendStatusToString 将int类型的SendStatus转为可读的中文描述
func SendStatusToString(status primitive.SendStatus) string {
	switch status {
	case primitive.SendOK:
		return "消息发送成功（已持久化）"
	case primitive.SendFlushDiskTimeout:
		return "消息发送成功但刷盘超时（仅存于Broker内存）"
	case primitive.SendFlushSlaveTimeout:
		return "消息发送成功但主从复制超时（从节点未同步）"
	case primitive.SendSlaveNotAvailable:
		return "消息发送成功但从节点不可用（仅主节点存储）"
	case primitive.SendUnknownError:
		return "消息发送未知错误"
	default:
		return fmt.Sprintf("未知发送状态(%d)", status)
	}
}
