package models

import (
	"time"
)

type UserRelation struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	UserId    uint       `json:"user_id"`
	FriendId  uint       `json:"friend_id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at"`
}

func (UserRelation) TableName() string {
	return "user_relation"
}
