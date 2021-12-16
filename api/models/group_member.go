package models

import (
	"gin-derived/global"
	"time"
)

type GroupMember struct {
	ID         uint       `gorm:"primary_key" json:"id"`
	GroupId    uint       `json:"group_id"`
	UserId     uint       `json:"user_id"`
	Username   string     `json:"username"`
	UserAvatar string     `json:"user_avatar"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `sql:"index" json:"deleted_at"`
}

func (GroupMember) TableName() string {
	return "group_member"
}

//FindGroups 查询自己所有群
func (gm *GroupMember) FindGroups() ([]*GroupMember, error) {
	var result []*GroupMember
	err := global.GDB.Select("group_id").Where(&gm).Find(&result).Error
	return result, err
}

func (gm *GroupMember) AddGroupMember(v string) (*GroupMember, error) {
	user := &User{}
	if err := global.GDB.Where(&User{Uuid: v}).Find(&user).Error; err != nil {
		return nil, err
	}
	gm.UserId = user.ID
	gm.Username = user.Username
	gm.UserAvatar = user.Avatar
	if err := global.GDB.Where("user_id = ?", gm.UserId).Where("group_id = ?", gm.GroupId).FirstOrCreate(&gm).Error; err != nil {
		return nil, err
	}
	return gm, nil
}
