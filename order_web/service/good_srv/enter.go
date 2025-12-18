package good_srv

import "order_web/common"

type GoodListRequest struct {
	common.PageInfo
	IsHot         bool  `form:"is_hot"`
	IsNew         bool  `form:"is_new"`
	PriceMax      int32 `form:"price_max"`
	PriceMin      int32 `form:"price_min"`
	BrandID       int32 `form:"brand_id"`
	TopCategoryID int32 `form:"top_category_id"`
}

type GoodCreateRequest struct {
	Name        string   `form:"name" json:"name" binding:"required,min=2,max=100"`
	GoodsSn     string   `form:"goods_sn" json:"goods_sn" binding:"required,min=2,lt=20"`
	Stocks      int32    `form:"stocks" json:"stocks" binding:"required,min=1"`
	CategoryId  int32    `form:"category" json:"category" binding:"required"`
	MarketPrice float32  `form:"market_price" json:"market_price" binding:"required"`
	ShopPrice   float32  `form:"shop_price" json:"shop_price" binding:"required,min=0"`
	GoodsBrief  string   `form:"goods_brief" json:"goods_brief" binding:"required,min=3"`
	Images      []string `form:"images" json:"images" binding:"required,min=1"`
	DescImages  []string `form:"desc_images" json:"desc_images" binding:"required"`
	ShipFree    *bool    `form:"ship_free" json:"ship_free" binding:"required"`
	FrontImage  string   `form:"front_image" json:"front_image" binding:"required,url"`
	Brand       int32    `form:"brand" json:"brand" binding:"required"`
}

type GoodDetailRequest struct {
	Id int32 `uri:"id" binding:"required,min=1"`
}

type GoodUpdateRequest struct {
	Name        string   `form:"name" json:"name" binding:"omitempty,min=2,max=100"`       // 有值时校验长度 2-100，
	GoodsSn     string   `form:"goods_sn" json:"goods_sn" binding:"omitempty,min=2,lt=20"` // 有值时校验长度 2-19，
	Stocks      int32    `form:"stocks" json:"stocks" binding:"omitempty,min=1"`           // 有值时校验 ≥1
	CategoryId  int32    `form:"category" json:"category"`
	MarketPrice float32  `form:"market_price" json:"market_price"`
	ShopPrice   float32  `form:"shop_price" json:"shop_price" binding:"omitempty,min=0"`
	GoodsBrief  string   `form:"goods_brief" json:"goods_brief" binding:"omitempty,min=3"`
	Images      []string `form:"images" json:"images" binding:"omitempty,min=1"`
	DescImages  []string `form:"desc_images" json:"desc_images"`
	ShipFree    *bool    `form:"ship_free" json:"ship_free"`
	FrontImage  string   `form:"front_image" json:"front_image"`
	Brand       int32    `form:"brand" json:"brand"`
}

type GoodPatchUpdateRequest struct {
	IsNew  *bool `form:"is_new" json:"is_new"`
	IsHot  *bool `form:"is_hot" json:"is_hot"`
	OnSale *bool `form:"on_sale" json:"on_sale"`
}

type GoodDeleteRequest struct {
	Id int32 `uri:"id" binding:"required,min=1"`
}
