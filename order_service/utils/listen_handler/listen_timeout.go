package listen_handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"go.uber.org/zap"
	"order_service/handler"
	"order_service/service"

	"order_service/global"
	"order_service/models"
)

func ListenMq() {
	messgaes, err := rocketmq.NewPushConsumer(consumer.WithNameServer([]string{global.Config.RocketMQ.Addr()}),
		consumer.WithGroupName(global.Config.RocketMQ.GroupName))
	if err != nil {
		panic(err)
	}
	messgaes.Subscribe(global.Config.RocketMQ.ConsumerSubscribe, consumer.MessageSelector{}, Timeout)
	messgaes.Start()
	select {}

}

func Timeout(ctx context.Context, msg ...*primitive.MessageExt) (consumer.ConsumeResult, error) {

	for i := range msg {
		var orderInfo service.OrderTransitionRequest
		_ = json.Unmarshal(msg[i].Body, &orderInfo)
		// 需要查看 订单存不存咋
		var orderModel models.OrderModel
		err := global.DB.Where(models.OrderModel{OrderSn: orderInfo.OrderSns}).Take(&orderModel).Error
		if err != nil {
			// 没找到 说明没有 这样就不需要管了 因为都没创建订单 所以不需要归还库存
			return consumer.ConsumeSuccess, nil
		}
		// 找到了  查一下 看看是不是已经支付了
		tx := global.DB.Begin()
		if orderModel.Status != "TRADE_SUCCESS" {
			// 说明没支付  我们要关闭， 然后发送消息给mq 让库存服务归还
			orderModel.Status = "CLOSED"
			err := tx.Save(&orderModel).Error
			if err != nil {
				zap.S().Error(err)
				tx.Rollback()
				return consumer.ConsumeRetryLater, nil
			}

			_, err = handler.GlobalOrderServer.MessageProducer.SendSync(context.Background(),
				primitive.NewMessage(global.Config.RocketMQ.ConsumerTopic, msg[i].Body))

			if err != nil {
				tx.Rollback()
				return consumer.ConsumeRetryLater, nil
			}
		}
		tx.Commit()

	}
	return consumer.ConsumeSuccess, nil
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
