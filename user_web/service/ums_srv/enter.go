package ums_srv

type RoleIdRequest struct {
	Id int32 `uri:"id" binding:"required,min=1"`
}

type RoleListRequest struct {
	Page  int32  `form:"page"`
	Limit int32  `form:"limit"`
	Key   string `form:"key"`
}

type RoleUpsertRequest struct {
	RoleName string  `json:"role_name" binding:"required"`
	MenuIds  []int32 `json:"menu_ids"`
}
