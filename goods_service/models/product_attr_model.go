package models

// ProductAttr 商品屬性表
type ProductAttr struct {
	Model
	AttrName  string `gorm:"type:varchar(100);not null;comment:屬性名稱;index:idx_product_attr_name"`
	AttrValue string `gorm:"type:varchar(255);not null;comment:屬性值"`
}
