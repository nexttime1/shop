package mq

import (
	"context"
	"encoding/json"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
	"time"

	"stock_service/global"
	"stock_service/models"
	"stock_service/proto"
	"stock_service/service"
)

func ListenMq() {
	messgaes, err := rocketmq.NewPushConsumer(consumer.WithNameServer([]string{global.Config.RocketMQ.Addr()}),
		consumer.WithGroupName(global.Config.RocketMQ.ConsumerGroupName))
	if err != nil {
		panic(err)
	}
	messgaes.Subscribe(global.Config.RocketMQ.ConsumerSubscribe, consumer.MessageSelector{}, AutoReBack)
	messgaes.Start()
	select {}

}

func AutoReBack(ctx context.Context, msg ...*primitive.MessageExt) (consumer.ConsumeResult, error) {

	for i := range msg {
		var orderInfo service.OrderTransitionRequest
		_ = json.Unmarshal(msg[i].Body, &orderInfo)
		// 开启事务
		tx := global.DB.Begin()
		rollbackFlag := true // 标记是否需要回滚
		defer func() {
			if r := recover(); r != nil {
				zap.S().Errorf("处理消息panic，OrderSn: %s, err: %v", orderInfo.OrderSns, r)
				tx.Rollback()
			} else if rollbackFlag {
				// 未主动提交时，回滚事务
				tx.Rollback()
			}
		}()
		//乐观锁
		// 需要查看 OrderSns 和 状态为1 也就是扣了为归还的 这个时候用了 乐观锁
		var history models.StockSellDetail
		err := tx.Where(models.StockSellDetail{OrderSn: orderInfo.OrderSns, Status: 1}).Take(&history).Error
		if err != nil {
			// 没找到 说明没有  或者 防止重复消费
			return consumer.ConsumeSuccess, nil
		}
		// 找到了  进行归还库存 并且 改历史记录 状态 变成2
		// 构造 Reback 函数需要的参数
		var info proto.SellInfo
		var list []*proto.GoodsInvInfo
		for _, inv := range history.Detail {
			list = append(list, &proto.GoodsInvInfo{
				GoodsId: inv.GoodId,
				Num:     inv.Num,
			})
		}
		info.GoodsInfo = list
		info.OrderSn = orderInfo.OrderSns
		_, err = Reback(tx, &info)
		if err != nil {
			tx.Rollback()
			return consumer.ConsumeRetryLater, err
		}
		// 改历史记录 状态 变成2  并且要必须是 我们拿到的版本 如果不是 正常有个for循环 咱这个环境下 我们就没必要for了 因为一会就延迟再来一次
		err = tx.Model(&history).
			Where("id = ? and version = ?", history.ID, history.Version). // 基于 ID 和 version 做乐观锁
			Updates(map[string]interface{}{
				"status":  2,
				"version": history.Version + 1,
			}).Error
		if err != nil {
			tx.Rollback()
			return consumer.ConsumeRetryLater, err
		}
		err = tx.Commit().Error
		if err != nil {
			tx.Rollback()
			return consumer.ConsumeRetryLater, err
		}
		rollbackFlag = false // 标记是否需要回滚

	}
	return consumer.ConsumeSuccess, nil
}

func Reback(tx *gorm.DB, info *proto.SellInfo) (*emptypb.Empty, error) {
	// 传递事务，保证操作原子性
	for _, invInfo := range info.GoodsInfo {
		// 乐观锁保证 高并发情况下 不会发生错误  比如两个请求同一个商品进行归还  读取的值都是100 都加50  防止最后是150
		retryCount := 0
		maxOptimisticRetry := 10
		for retryCount < maxOptimisticRetry {
			var model models.InventoryModel
			err := tx.Where("goods = ?", invInfo.GoodsId).Take(&model).Error
			if err != nil {
				zap.S().Error(err)
				return nil, status.Error(codes.NotFound, "商品库存不存在")
			}
			// 库存 +
			model.Stock += invInfo.Num

			err = tx.Model(models.InventoryModel{}).Where("goods = ? and version = ?", model.Goods, model.Version).Select("stock", "version").Updates(map[string]interface{}{"stock": model.Stock, "version": model.Version + 1}).Error
			if err != nil {
				retryCount++
				zap.S().Warnf("商品%d乐观锁重试，当前次数: %d", invInfo.GoodsId, retryCount)
				time.Sleep(100 * time.Millisecond)
				continue
			} else {
				break
			}
		}
		// 重试次数耗尽仍未成功
		if retryCount >= maxOptimisticRetry {
			zap.S().Errorf("商品%d乐观锁重试次数耗尽，更新失败", invInfo.GoodsId)
			return nil, status.Error(codes.Internal, "库存更新并发冲突，请重试")
		}
	}
	return &emptypb.Empty{}, nil
}
