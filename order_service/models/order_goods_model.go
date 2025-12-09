package models

// OrderGoodsModel 订单商品明细表模型
type OrderGoodsModel struct {
	Model
	Order      int32   `gorm:"type:int;index;comment:订单ID"`
	Goods      int32   `gorm:"type:int;index;comment:商品ID"`
	GoodsName  string  `gorm:"type:varchar(100);comment:商品名称"`
	GoodsPrice float32 `gorm:"comment:商品单价"`
	GoodImages string  `gorm:"type:varchar(100);comment:商品图片"`
	Nums       int32   `gorm:"type:int;comment:商品数量"`
}

// TableName 重写订单商品明细表名
func (OrderGoodsModel) TableName() string {
	return "ordergoods"
}
