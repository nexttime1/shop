package brand_srv

type BrandCreateRequest struct {
	Name string `form:"name" json:"name" binding:"required,min=3,max=10"`
	Logo string `form:"logo" json:"logo" binding:"url"`
}

type BrandUpdateRequest struct {
	Name string `form:"name" json:"name"`
	Logo string `form:"logo" json:"logo"`
}

type BrandIdRequest struct {
	Id int32 `uri:"id" binding:"required,min=1"`
}

type CreateCategoryBrandRequest struct {
	CategoryId int32 `json:"category_id" binding:"required"`
	BrandId    int32 `json:"brand_id" binding:"required"`
}

type UpdateCategoryBrandRequest struct {
	CategoryId int32 `json:"category_id" `
	BrandId    int32 `json:"brand_id" `
}
