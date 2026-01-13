package pms_srv

type ProductAttrListRequest struct {
	Page  int32  `form:"page"`
	Limit int32  `form:"limit"`
	Key   string `form:"key"`
}

type ProductAttrCreateRequest struct {
	AttrName  string `json:"attr_name" binding:"required"`
	AttrValue string `json:"attr_value" binding:"required"`
}

type SkuStockListRequest struct {
	Page      int32 `form:"page"`
	Limit     int32 `form:"limit"`
	ProductID int32 `form:"product_id"`
}

type SkuStockCreateRequest struct {
	ProductID int32  `json:"product_id" binding:"required"`
	SkuCode   string `json:"sku_code" binding:"required"`
	Stock     int32  `json:"stock" binding:"required"`
}
