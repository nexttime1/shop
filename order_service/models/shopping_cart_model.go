package models

// ShoppingCartModel
type ShoppingCartModel struct {
	Model
	User    int32 `gorm:"type:int;index;comment:用户ID"`
	Goods   int32 `gorm:"type:int;index;comment:商品ID"`
	Nums    int32 `gorm:"type:int;comment:商品数量"`
	Checked *bool `gorm:"comment:是否勾选（结算）"`
}

// TableName 重写购物车表名
func (ShoppingCartModel) TableName() string {
	return "shoppingcart"
}
