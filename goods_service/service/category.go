package service

type CategoryUpdateServiceMap struct {
	Name             string `structs:"name"`
	ParentCategoryID int32  `structs:"parent_category_id"`
	Level            int32  `structs:"level"`
	IsTab            *bool  `structs:"is_tab"`
}

type CategoryBrandUpdateServiceMap struct {
	CategoryId int32 `structs:"category_id"`
	BrandId    int32 `structs:"brands_id"`
}
