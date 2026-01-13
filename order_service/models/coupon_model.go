package models

// Coupon 優惠券
type Coupon struct {
	Model
	CouponCode string `gorm:"type:varchar(64);not null;uniqueIndex:idx_coupon_code;comment:優惠券碼"`
	Title      string `gorm:"type:varchar(128);not null;comment:標題"`
	Discount   int    `gorm:"type:int;not null;default:0;comment:折扣值"`
}
