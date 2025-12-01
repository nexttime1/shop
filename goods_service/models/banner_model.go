package models

type BannerModel struct {
	Model `structs:"-"`
	Image string `gorm:"type:varchar(200);not null;comment: 图片的url" structs:"image"`
	Url   string `gorm:"type:varchar(200);not null;comment:跳转的详情"  structs:"url"`
	Index int32  `gorm:"type:int;default:1;not null"  structs:"index"`
}
