package handler

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"stock_service/global"
	"stock_service/models"
	"stock_service/proto"
	"time"
)

type InventorySever struct {
}

func (i InventorySever) SetInv(ctx context.Context, info *proto.GoodsInvInfo) (*emptypb.Empty, error) {
	// 有就设置 没有就新增
	var model models.InventoryModel
	count := global.DB.Where("goods = ?", info.GoodsId).First(&model).RowsAffected
	if count == 0 {
		// 新增
		err := global.DB.Create(&models.InventoryModel{
			Goods: info.GoodsId,
			Stock: info.Num,
		}).Error
		if err != nil {
			zap.S().Error(err)
			return nil, status.Errorf(codes.Internal, "创建失败")
		}
		return &emptypb.Empty{}, nil
	}
	model.Stock = info.Num
	err := global.DB.Save(&model).Error
	if err != nil {
		zap.S().Error(err)
		return nil, status.Errorf(codes.Internal, "更新失败")
	}
	return &emptypb.Empty{}, nil

}

func (i InventorySever) InvDetail(ctx context.Context, info *proto.GoodsInvInfo) (*proto.GoodsInvInfo, error) {
	var model models.InventoryModel
	err := global.DB.Where("goods = ?", info.GoodsId).First(&model).Error
	if err != nil {
		zap.S().Error(err)
		return nil, status.Errorf(codes.NotFound, "库存信息不存在")
	}
	return &proto.GoodsInvInfo{
		GoodsId: info.GoodsId,
		Num:     model.Stock,
	}, nil

}

func (i InventorySever) Sell(ctx context.Context, info *proto.SellInfo) (*emptypb.Empty, error) {
	// 开启事务，保证操作原子性
	tx := global.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 生成历史记录
	var record models.StockSellDetail
	record.OrderSn = info.OrderSn
	record.Status = 1 // 扣减未归还

	var DetailList []models.GoodsDetail
	for _, invInfo := range info.GoodsInfo {
		DetailList = append(DetailList, models.GoodsDetail{
			GoodId: invInfo.GoodsId,
			Num:    invInfo.Num,
		})
		// redis 分布式锁
		mutexName := fmt.Sprintf("good_%d", invInfo.GoodsId)
		mutex := global.RedisMutex.NewMutex(mutexName)
		//for {	// 乐观锁  也先不用了
		var model models.InventoryModel
		// 悲观锁
		//err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("goods = ?", invInfo.GoodsId).Take(&model).Error
		err := mutex.Lock()
		if err != nil {
			zap.S().Error(err)
			return nil, status.Errorf(codes.Internal, "redis分布式锁加载错误")
		}
		err = tx.Where("goods = ?", invInfo.GoodsId).Take(&model).Error
		if err != nil {
			zap.S().Error(err)
			tx.Rollback()
			return nil, status.Error(codes.NotFound, "商品库存不存在")
		}
		if model.Stock < invInfo.Num {
			zap.S().Error(err)
			tx.Rollback()
			return nil, status.Error(codes.ResourceExhausted, "库存不足")
		}
		// 库存 -
		model.Stock -= invInfo.Num
		err = tx.Save(&model).Error
		if err != nil {
			zap.S().Error(err)
			tx.Rollback()
			return nil, status.Error(codes.Internal, "更新错误")
		}
		ok, err := mutex.Unlock()
		if err != nil || !ok {
			tx.Rollback()
			zap.S().Error(err)
			return nil, status.Error(codes.Internal, "解锁失败")
		}

		//err = tx.Model(models.InventoryModel{}).Where("goods = ? and version = ?", model.Goods, model.Version).Select("stock", "version").Updates(map[string]interface{}{"stock": model.Stock, "version": model.Version + 1}).Error
		//if err != nil {
		//	continue
		//} else {
		//	break
		//}
		//}
	}
	record.Detail = DetailList
	err := tx.Create(&record).Error
	if err != nil {
		tx.Rollback()
		return nil, status.Error(codes.Internal, "创建订单历史记录错误")
	}

	return &emptypb.Empty{}, tx.Commit().Error

}

func (i InventorySever) Reback(ctx context.Context, info *proto.SellInfo) (*emptypb.Empty, error) {
	// 开启事务，保证操作原子性
	tx := global.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	for _, invInfo := range info.GoodsInfo {
		// 乐观锁保证 高并发情况下 不会发生错误  比如两个请求同一个商品进行归还  读取的值都是100 都加50  防止最后是150
		retryCount := 0
		maxOptimisticRetry := 10
		for retryCount < maxOptimisticRetry {
			var model models.InventoryModel
			err := tx.Where("goods = ?", invInfo.GoodsId).Take(&model).Error
			if err != nil {
				zap.S().Error(err)
				tx.Rollback()
				return nil, status.Error(codes.NotFound, "商品库存不存在")
			}
			// 库存 +
			model.Stock += invInfo.Num

			err = tx.Model(models.InventoryModel{}).Where("goods = ? and version = ?", model.Goods, model.Version).Select("stock", "version").Updates(map[string]interface{}{"stock": model.Stock, "version": model.Version + 1}).Error
			if err != nil {
				retryCount++
				zap.S().Warnf("商品%d乐观锁重试，当前次数: %d", invInfo.GoodsId, retryCount)
				time.Sleep(10 * time.Millisecond)
				continue
			} else {
				break
			}
		}
		// 重试次数耗尽仍未成功
		if retryCount >= maxOptimisticRetry {
			_ = tx.Rollback()
			zap.S().Errorf("商品%d乐观锁重试次数耗尽，更新失败", invInfo.GoodsId)
			return nil, status.Error(codes.Internal, "库存更新并发冲突，请重试")
		}
	}
	return &emptypb.Empty{}, tx.Commit().Error
}
