package models

type UserCollectionModel struct {
	Model
	UserId int32 `gorm:"type:int;index:idx_user_goods,unique"`
	GoodId int32 `gorm:"type:int;index:idx_user_goods,unique"`
}
