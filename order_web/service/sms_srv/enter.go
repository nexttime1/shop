package sms_srv

type CouponListRequest struct {
	Page  int32 `form:"page"`
	Limit int32 `form:"limit"`
}

type CouponCreateRequest struct {
	CouponCode string `json:"coupon_code" binding:"required"`
	Title      string `json:"title" binding:"required"`
	Discount   int32  `json:"discount" binding:"required"`
}

type FlashCreateRequest struct {
	Name    string `json:"name" binding:"required"`
	FlashID int32  `json:"flash_id" binding:"required"`
}

type AdCreateRequest struct {
	Image string `json:"image" binding:"required,url"`
	Url   string `json:"url" binding:"required"`
}
