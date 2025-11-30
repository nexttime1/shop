package models

type GoodModel struct {
	Model
	CategoryID int32          `gorm:"type:int;not null;comment:分类ID（逻辑外键）;index:idx_goods_category"`
	Category   *CategoryModel `gorm:"foreignKey:CategoryID;references:ID;constraint:<-:false,foreignKey:no action"`

	BrandsID int32   `gorm:"type:int;not null;comment:品牌ID（逻辑外键）;index:idx_goods_brand"`
	Brands   *Brands `gorm:"foreignKey:BrandsID;references:ID;constraint:<-:false,foreignKey:no action"`

	OnSale      bool    `gorm:"default:false;not null;comment:是否上架"`
	ShipFree    bool    `gorm:"default:false;not null;comment:是否包邮"`
	IsNew       bool    `gorm:"default:false;not null;comment:是否新品"`
	IsHot       bool    `gorm:"default:false;not null;comment:是否热销"`
	Name        string  `gorm:"type:varchar(50);not null;comment:商品名称;index:idx_goods_name"`
	GoodsSn     string  `gorm:"type:varchar(50);not null;comment:商品编号;uniqueIndex:idx_goods_sn"`
	ClickNum    int32   `gorm:"type:int;default:0;not null;comment:点击量"`
	SoldNum     int32   `gorm:"type:int;default:0;not null;comment:销量"`
	FavNum      int32   `gorm:"type:int;default:0;not null;comment:收藏量"`
	MarketPrice float32 `gorm:"not null;comment:市场价"`
	ShopPrice   float32 `gorm:"not null;comment:售价;index:idx_goods_price"`
	GoodsBrief  string  `gorm:"type:varchar(100);not null;comment:商品简介"`

	// 方便查询商品的所有图片（Gorm虚拟字段，不存数据库）
	Images []*GoodsImageModel `gorm:"foreignKey:GoodsID;references:ID;constraint:<-:false,foreignKey:no action"`
}
