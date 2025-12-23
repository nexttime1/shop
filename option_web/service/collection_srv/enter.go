package collection_srv

type CollectionRequest struct {
	GoodID int32 `json:"good_id" binding:"required"`
}

type CartAddResponse struct {
	Id int32 `json:"id"`
}

type CartIdRequest struct {
	Id int32 `uri:"id" binding:"required,min=1"`
}

type CartUpdateRequest struct {
	Num     int32 `json:"num" binding:"required" min:"1"`
	Checked *bool `json:"checked"`
}
