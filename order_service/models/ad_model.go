package models

// Ad 廣告
type Ad struct {
	Model
	Image string `gorm:"type:varchar(255);not null;comment:圖片地址"`
	Url   string `gorm:"type:varchar(255);not null;comment:跳轉鏈接"`
}
