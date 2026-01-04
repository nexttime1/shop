package mq

import (
	"context"
	"encoding/json"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"go.uber.org/zap"

	"order_service/global"
	"order_service/models"
	"order_service/service"
)

func ListenMq() {
	messgaes, err := rocketmq.NewPushConsumer(consumer.WithNameServer([]string{"192.168.163.132:9876"}),
		consumer.WithGroupName("shop_order"))
	if err != nil {
		panic(err)
	}
	messgaes.Subscribe("order_timeout", consumer.MessageSelector{}, Timeout)
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
			p, err := rocketmq.NewProducer(producer.WithNameServer([]string{"192.168.163.132:9876"}))
			if err != nil {
				tx.Rollback()
				return consumer.ConsumeRetryLater, nil
			}
			err = p.Start()
			if err != nil {
				tx.Rollback()
				return consumer.ConsumeRetryLater, nil
			}

			_, err = p.SendSync(context.Background(), primitive.NewMessage("shop_reback", msg[i].Body))

			if err != nil {
				tx.Rollback()
				return consumer.ConsumeRetryLater, nil
			}
		}
		tx.Commit()

	}
	return consumer.ConsumeSuccess, nil
}
