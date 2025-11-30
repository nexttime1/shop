package models

type CategoryModel struct {
	Model
	Name             string `gorm:"type:varchar(20);not null"`
	ParentCategoryID int32  `gorm:"comment:父分类ID（逻辑外键，关联自身ID）"`
	// 用 constraint 禁用物理外键约束   保留是因为 查询方便
	ParentCategory *CategoryModel `gorm:"foreignKey:ParentCategoryID;references:ID;constraint:<-:false,foreignKey:no action"`
	Level          int32          `gorm:"type:int;not null;default:1"`
	IsTab          bool           `gorm:"default:false;not null"`
}
