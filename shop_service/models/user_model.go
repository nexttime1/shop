package models

import "time"

type UserModel struct {
	Model    `structs:"-"`
	Mobile   string     `gorm:"index:idx_mobile;unique;type:varchar(11);not null" structs:"-"`
	Password string     `gorm:"type:varchar(100);not null" structs:"password"`
	NickName string     `gorm:"type:varchar(100);"  structs:"nick_name"`
	Birthday *time.Time `gorm:"type:datetime" structs:"birthday"`
	Gender   string     `gorm:"column:gender;default:male;type:varchar(6)"  structs:"gender"`
	Role     int        `gorm:"column: role;default 1"  structs:"role"`
}
