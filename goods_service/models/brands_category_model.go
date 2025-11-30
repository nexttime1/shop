package models

type BrandCategoryModel struct {
	Model
	CategoryID int32 `gorm:"type:int;not null;comment:分类ID（逻辑外键，关联category_models.id）;index:idx_category_brand,unique"`
	//禁用物理外键约束
	Category *CategoryModel `gorm:"foreignKey:CategoryID;references:ID;constraint:<-:false,foreignKey:no action"`
	BrandsID int32          `gorm:"type:int;not null;comment:品牌ID（逻辑外键，关联brands.id）;index:idx_category_brand,unique"`
	// 禁用物理外键约束
	Brands *Brands `gorm:"foreignKey:BrandsID;references:ID;constraint:<-:false,foreignKey:no action"`
}
