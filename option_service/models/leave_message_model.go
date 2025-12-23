package models

import "option_service/models/enum"

type LeavingMessageModel struct {
	Model
	UserId      int32            `gorm:"type:int(11);index"`
	MessageType enum.MessageType `gorm:"type:int(11)"`
	Subject     string           `gorm:"type:varchar(128)"`
	Message     string           `gorm:"type:varchar(128)"`
	File        string           `gorm:"type:varchar(200)"`
}
