package mq

import (
	"context"
	"encoding/json"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"gorm.io/gorm"
	"stock_service/global"
	"stock_service/models"
	"stock_service/service"
)

func ListenMq() {
	messgaes, err := rocketmq.NewPushConsumer(consumer.WithNameServer([]string{"192.168.163.132:9876"}),
		consumer.WithGroupName("shop_inventory"))
	if err != nil {
		panic(err)
	}
	messgaes.Start()
	defer messgaes.Shutdown()
	messgaes.Subscribe("shop_reback", consumer.MessageSelector{}, AutoReBack)
}

func AutoReBack(ctx context.Context, msg ...*primitive.MessageExt) (consumer.ConsumeResult, error) {

	for i := range msg {
		var orderInfo service.OrderTransitionRequest
		_ = json.Unmarshal(msg[i].Body, &orderInfo)
		// 需要查看 OrderSns 和 状态为1 也就是扣了为归还的
		var history models.StockSellDetail
		err := global.DB.Where(models.StockSellDetail{OrderSn: orderInfo.OrderSns, Status: 1}).Take(&history).Error
		if err != nil {
			// 没找到 说明没有
			return consumer.ConsumeSuccess, nil
		}
		// 找到了  进行归还库存 并且 改历史记录 状态 变成2
		tx := global.DB.Begin()
		for _, inv := range history.Detail {
			err := tx.Model(models.InventoryModel{Goods: inv.GoodId}).Update("stock", gorm.Expr("stock + ?", inv.Num)).Error
			if err != nil {
				tx.Rollback()
				return consumer.ConsumeRetryLater, err
			}
		}
		err = tx.Model(&history).Update("status", 1).Error
		if err != nil {
			tx.Rollback()
			return consumer.ConsumeRetryLater, err
		}

	}
	return consumer.ConsumeSuccess, nil
}
