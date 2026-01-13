package models

import "gorm.io/datatypes"

// Role 權限角色表
type Role struct {
	Model
	RoleName string         `gorm:"type:varchar(64);not null;uniqueIndex:idx_role_name;comment:角色名稱"`
	MenuIDs  datatypes.JSON `gorm:"type:json;not null;comment:菜單ID集合(JSON)"`
}
