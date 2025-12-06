package service

type GoodUpdateServiceMap struct {
	Name        string  `structs:"name"`
	GoodsSn     string  `structs:"goods_sn"`
	Stocks      int32   `structs:"stocks"`
	MarketPrice float32 `structs:"market_price"`
	ShopPrice   float32 `structs:"shop_price"`
	GoodsBrief  string  `structs:"goods_brief"`
	ShipFree    *bool   `structs:"ship_free"`
	IsNew       *bool   `structs:"is_new"`
	IsHot       *bool   `structs:"is_hot"`
	OnSale      *bool   `structs:"on_sale"`
	CategoryId  int32   `structs:"category_id"`
	Brand       int32   `structs:"brands_id"`
}
