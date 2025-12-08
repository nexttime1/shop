package models

type InventoryModel struct {
	Model   `structs:"-"`
	Goods   int32 `gorm:"type:int;index"`
	Stock   int32 `gorm:"type:int"`
	Version int32 `gorm:"type:int"` //分布式锁
}
