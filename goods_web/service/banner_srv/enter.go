package banner_srv

type BannerCreateRequest struct {
	Image string `form:"image" json:"image" binding:"url"`
	Index int32  `form:"index" json:"index" binding:"required"`
	Url   string `form:"url" json:"url" binding:"url"`
}

type BannerIdRequest struct {
	Id int32 `uri:"id" binding:"required,min=1"`
}

type BannerUpdateRequest struct {
	Image string `form:"image" json:"image"`
	Index int32  `form:"index" json:"index"`
	Url   string `form:"url" json:"url"`
}
