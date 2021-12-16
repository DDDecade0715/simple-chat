package models

import (
	"gin-derived/global"
	"time"
)

type ChatImage struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	MessageId string     `json:"message_id"`
	Url       string     `json:"url"`
	FileName  string     `json:"file_name"`
	FileSize  int64      `json:"file_size"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at"`
}

func (ChatImage) TableName() string {
	return "chat_images"
}

func CreateChatImage(value *ChatImage) (err error) {
	err = global.GDB.Create(&value).Error
	return
}

func FindChatImage(value *ChatImage) (result *ChatImage, err error) {
	result = &ChatImage{}
	err = global.GDB.Where(&value).First(&result).Error
	return
}
