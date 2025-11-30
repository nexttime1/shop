package models

type Brands struct {
	Model
	Name string `gorm:"type:varchar(20);not null"`
	Logo string `gorm:"type:varchar(200);default:'';not null"`
}
