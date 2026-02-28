package listen_handler

import (
	"context"
	"encoding/json"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"go.uber.org/zap"
	"order_service/handler"
	"order_service/service"
	"time"

	"order_service/global"
	"order_service/models"
)

func ListenMq() {
	messgaes, err := rocketmq.NewPushConsumer(consumer.WithNameServer([]string{global.Config.RocketMQ.Addr()}),
		consumer.WithGroupName(global.Config.RocketMQ.TimeOutConsumerGroupName),
		// 最大重试次数
		consumer.WithMaxReconsumeTimes(global.Config.RocketMQ.MaxRetryTimes),
		// 重试延迟时间
		consumer.WithSuspendCurrentQueueTimeMillis(time.Duration(global.Config.RocketMQ.BaseRetryDelay)),
	)
	if err != nil {
		panic(err)
	}
	// 监听普通消息 订单时间超时
	messgaes.Subscribe(global.Config.RocketMQ.TimeOutTopic, consumer.MessageSelector{}, Timeout)
	messgaes.Start()
	select {}

}

func Timeout(ctx context.Context, msg ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	zap.S().Info("开始处理订单超时消息，本次处理消息数：", len(msg))
	// 标记是否有消息需要重试（整体返回结果）
	needRetry := false

	for i := range msg {
		// 单条消息独立的日志字段（便于定位问题）
		logFields := []zap.Field{
			zap.Int("index", i),
			zap.String("msgId", msg[i].MsgId),
			zap.String("orderSn", ""), // 后续解析后补充
		}

		// 解析单条消息体（独立错误处理，不影响其他消息）
		var orderInfo service.OrderTransitionRequest
		if err := json.Unmarshal(msg[i].Body, &orderInfo); err != nil {
			zap.S().Error("解析订单超时消息体失败", append(logFields, zap.Error(err)))
			// 解析失败：标记无需重试（消息本身有问题，重试也没用），继续处理下一条
			continue
		}
		logFields = append(logFields, zap.String("orderSn", orderInfo.OrderSns))

		// 2. 校验订单是否存在
		var orderModel models.OrderModel
		err := global.DB.Where(models.OrderModel{OrderSn: orderInfo.OrderSns}).Take(&orderModel).Error
		if err != nil {
			zap.S().Warn("超时订单不存在，无需处理", append(logFields, zap.Error(err)))
			// 订单不存在：继续处理下一条，不返回、不重试
			continue
		}

		// 3. 处理订单状态
		tx := global.DB.Begin()
		// 标记当前消息是否处理成功
		success := false
		defer func() {
			if !success {
				tx.Rollback()
			}
		}()

		if orderModel.Status != "TRADE_SUCCESS" {
			// 3.1 更新订单为已关闭
			orderModel.Status = "CLOSED"
			if err := tx.Save(&orderModel).Error; err != nil {
				zap.S().Error("更新订单状态为CLOSED失败", append(logFields, zap.Error(err)))
				needRetry = true // 标记需要整体重试
				continue         // 继续处理下一条
			}

			// 3.2 发送库存归还消息
			_, err := handler.GlobalOrderServer.MessageProducer.SendSync(ctx,
				primitive.NewMessage(global.Config.RocketMQ.StockTimeoutTopic, msg[i].Body))
			if err != nil {
				zap.S().Error("发送库存归还消息失败", append(logFields, zap.Error(err)))
				needRetry = true // 标记需要整体重试
				continue         // 继续处理下一条
			}
		}

		// 4. 提交事务（当前消息处理成功）
		if err := tx.Commit().Error; err != nil {
			zap.S().Error("提交订单超时事务失败", append(logFields, zap.Error(err)))
			needRetry = true
			continue
		}
		success = true
		zap.S().Info("订单超时处理成功", logFields)
	}

	// 根据是否有失败消息，返回整体消费结果
	if needRetry {
		// 有消息处理失败，返回重试（客户端会重新推送所有未处理成功的消息）
		return consumer.ConsumeRetryLater, nil
	}
	// 所有消息处理完成（成功/无需处理），返回成功
	return consumer.ConsumeSuccess, nil
}
