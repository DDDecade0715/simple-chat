package models

import (
	"gin-derived/global"
	"time"
)

type Group struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	UserId    uint       `json:"user_id"`
	Name      string     `json:"name"`
	Uuid      string     `sql:"index" json:"uuid"`
	Avatar    string     `json:"avatar"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at"`
}

func (Group) TableName() string {
	return "group"
}

func (g *Group) CreateGroup(contactIds []string) (err error) {
	tx := global.GDB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if tx.Error != nil {
		return err
	}
	if err = tx.Create(&g).Error; err != nil {
		tx.Rollback()
		return err
	}
	for _, v := range contactIds {
		user := &User{}
		if err = tx.Where(&User{Uuid: v}).Find(&user).Error; err != nil {
			tx.Rollback()
			return err
		}
		groupMember := &GroupMember{
			GroupId:    g.ID,
			UserId:     user.ID,
			Username:   user.Username,
			UserAvatar: user.Avatar,
		}
		if err = tx.Create(&groupMember).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

type GroupInfo struct {
	Id      uint               `json:"id"`
	Uuid    string             `json:"uuid"`
	Name    string             `json:"name"`
	Avatar  string             `json:"avatar"`
	Message *ContactsMessage   `json:"message"`
	Members []*ResponseMembers `json:"members"`
}

func (g *Group) FindGroupInfo(value []uint) (result []*GroupInfo, err error) {
	err = global.GDB.Model(&g).Where("id in (?)", value).Scan(&result).Error
	return
}

type ResponseMembers struct {
	Name       string `json:"name"`
	GroupUuid  string `json:"group_uuid"`
	GroupId    uint   `json:"group_id"`
	UserId     uint   `json:"user_id"`
	Username   string `json:"username"`
	UserAvatar string `json:"user_avatar"`
	Uuid       string `json:"uuid"`
}

func (g *Group) FindGroupData() (result []*ResponseMembers, err error) {
	err = global.GDB.Table("group").Select("group.name,group.uuid as group_uuid,group_member.group_id,group_member.user_id,group_member.username,group_member.user_avatar,users.uuid").Joins("left join group_member on group.id = group_member.group_id").Joins("left join users on group_member.user_id = users.id").Where(&g).Scan(&result).Error
	return
}

func (g *Group) FindGroup() (result *Group, err error) {
	result = &Group{}
	err = global.GDB.Where("uuid = ?", g.Uuid).First(&result).Error
	return
}

func (g *Group) UpdateGroup(update *Group) (*Group, error) {
	err := global.GDB.Model(&g).Where("uuid = ?", g.Uuid).Updates(&update).Error
	return g, err
}
