package good_srv

import "goods_api/common"

type GoodListRequest struct {
	common.PageInfo
	IsHot         bool  `form:"is_hot"`
	IsNew         bool  `form:"is_new"`
	PriceMax      int32 `form:"price_max"`
	PriceMin      int32 `form:"price_min"`
	BrandID       int32 `form:"brand_id"`
	TopCategoryID int32 `form:"top_category_id"`
}
