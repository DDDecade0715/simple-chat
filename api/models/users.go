package models

import (
	"gin-derived/global"
	"time"
)

type User struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	Username  string     `json:"username"`
	Password  string     `json:"password"`
	Avatar    string     `json:"avatar"`
	Uuid      string     `sql:"index" json:"uuid"`
	IsAdmin   uint       `json:"is_admin"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at"`
}

func (User) TableName() string {
	return "users"
}

func FindUser(value *User) (result *User, err error) {
	result = &User{}
	err = global.GDB.Where(value).First(&result).Error
	return result, err
}

func CreateUser(value *User) (result *User, err error) {
	result = &User{
		Username: value.Username,
		Password: value.Password,
		Uuid:     value.Uuid,
	}
	err = global.GDB.Create(&result).Error
	return
}

type Contacts struct {
	Uuid     string           `json:"uuid"`
	Username string           `json:"username"`
	Avatar   string           `json:"avatar"`
	Message  *ContactsMessage `json:"message"`
}

func FindUserContacts(value *User) (contacts []*Contacts, err error) {
	var result []*User
	err = global.GDB.Select("username,avatar,uuid").Not(value).Find(&result).Scan(&contacts).Error
	return
}

func UpdateUserinfo(value *User) (err error) {
	whereUser := &User{
		ID: value.ID,
	}
	updateUser := &User{
		Username: value.Username,
		Avatar:   value.Avatar,
	}
	err = global.GDB.Model(&whereUser).Updates(updateUser).Error
	return
}

func FindUsername(value *User) (result *User, err error) {
	result = &User{}
	err = global.GDB.Where("username = ?", value.Username).Where("id <> ?", value.ID).First(&result).Error
	return result, err
}
