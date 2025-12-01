package models

type Brands struct {
	Model `structs:"-"`
	Name  string `gorm:"type:varchar(20);not null" structs:"name"`
	Logo  string `gorm:"type:varchar(200);default:'';not null"  structs:"logo"`
}
