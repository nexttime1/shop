package models

// Flash 秒殺活動
type Flash struct {
	Model
	Name    string `gorm:"type:varchar(128);not null;comment:活動名稱"`
	FlashID int    `gorm:"type:int;not null;comment:活動標識"`
}
