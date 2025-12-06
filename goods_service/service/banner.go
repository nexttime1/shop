package service

type BannerUpdateServiceMap struct {
	Image string `structs:"image"`
	Url   string `structs:"url"`
	Index int32  `structs:"index"`
}
