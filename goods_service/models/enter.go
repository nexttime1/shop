package models

import (
	"gorm.io/gorm"
	"time"
)

type Model struct {
	ID        int32          `gorm:"primarykey" structs:"-" json:"id"`
	CreatedAt time.Time      `gorm:"column:add_time" structs:"-" json:"-"`
	UpdatedAt time.Time      `gorm:"column:update_time" structs:"-" json:"-"`
	DeletedAt gorm.DeletedAt `structs:"-" json:"-"`
	idDeleted bool           `structs:"-" json:"-"`
}
