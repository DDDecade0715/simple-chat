package models

import (
	"time"
)

type SystemMessage struct {
	ID          uint       `gorm:"primary_key" json:"id"`
	FromUserId  string     `json:"fromUserId"`
	ToContactId string     `json:"toContactId"`
	Type        string     `json:"type"`
	Content     string     `json:"content"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `sql:"index" json:"deleted_at"`
}

func (SystemMessage) TableName() string {
	return "system_message"
}
