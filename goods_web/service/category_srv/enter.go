package category_srv

type CategoryIdRequest struct {
	Id int32 `uri:"id" binding:"required,min=1"`
}

type CategoryCreateRequest struct {
	Name           string `form:"name" json:"name" binding:"required,min=3,max=20"`
	ParentCategory int32  `form:"parent" json:"parent"`
	Level          int32  `form:"level" json:"level" binding:"required,oneof=1 2 3"`
	IsTab          *bool  `form:"is_tab" json:"is_tab" binding:"required"`
}

type UpdateCategoryRequest struct {
	Name  string `form:"name" json:"name" binding:"required,min=3,max=20"`
	IsTab *bool  `form:"is_tab" json:"is_tab"`
}
