package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"go.uber.org/zap"
)

type InventoryModel struct {
	Model   `structs:"-"`
	Goods   int32 `gorm:"type:int;index"`
	Stock   int32 `gorm:"type:int"`
	Version int32 `gorm:"type:int"` //分布式锁
}

type StockSellDetail struct {
	OrderSn string          `gorm:"type:varchar(200);index:unique"`
	Status  int32           //1 表示已扣减 2. 表示已归还
	Detail  GoodsDetailList `gorm:"type:varchar(200)"`
}

type GoodsDetailList []GoodsDetail
type GoodsDetail struct {
	GoodId int32
	Num    int32
}

func (g *GoodsDetailList) Value() (driver.Value, error) {
	return json.Marshal(g)
}

func (g *GoodsDetailList) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		zap.S().Error("错误：type assertion to []byte failed")
		return errors.New("type assertion to []byte failed")
	}
	err := json.Unmarshal(b, g)
	return err
}
