package models

import (
	"fmt"
	"gin-derived/global"
	"github.com/jinzhu/gorm"
	"math"
	"time"
)

type Message struct {
	ID          uint       `gorm:"primary_key" json:"id"`
	MessageId   string     `json:"messageId"`
	FromUserId  string     `json:"fromUserId"`
	ToContactId string     `json:"toContactId"`
	Type        string     `json:"type"`
	Content     string     `json:"content"`
	FromUser    string     `json:"fromUser"`
	Status      string     `json:"status"`
	SendTime    string     `time_format:"2006-01-02 15:04:05" json:"sendTime"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `sql:"index" json:"deleted_at"`
}

func (Message) TableName() string {
	return "messages"
}

//MessageChat 聊天信息体
type MessageChat struct {
	Id          string    `json:"id"`
	Content     string    `json:"content"`
	FromUser    *FromUser `json:"fromUser"`
	SendTime    int64     `json:"sendTime"`
	Status      string    `json:"status"`
	ToContactId string    `json:"toContactId"`
	Type        string    `json:"type"`
	MessageId   string    `json:"messageId"`
	FileName    string    `json:"fileName"`
	FileSize    int64     `json:"fileSize"`
}

type FromUser struct {
	Avatar      string `json:"avatar"`
	DisplayName string `json:"displayName"`
	Id          string `json:"id"`
}

//CreateGroupMessage 创建群聊信息体
type CreateGroupMessage struct {
	Id      string   `json:"id"`
	Type    string   `json:"type"`
	Members []string `json:"members"`
}

//CreateGroupMessageResponse 创建群聊信息体返回
type CreateGroupMessageResponse struct {
	Type        string `json:"type"`
	ToContactId string `json:"toContactId"`
	*GroupInfo
}

type ContactsMessage struct {
	Type     string `json:"type"`
	Content  string `json:"content"`
	SendTime string `time_format:"2006-01-02 15:04:05" json:"send_time"`
	Status   string `json:"status"`
}

func GetContactsMessage(whereFrom *Message, whereContact *Message) (contactsMessage *ContactsMessage, err error) {
	result := &Message{}
	contactsMessage = &ContactsMessage{}
	err = global.GDB.Where(whereFrom).Or(whereContact).Order("send_time desc").First(&result).Scan(&contactsMessage).Error
	return
}

func SaveOrUpdateMessage(value *Message) (err error) {
	if res, err := FindMessage(value); err != nil {
		//创建消息
		if gorm.IsRecordNotFoundError(err) {
			err = global.GDB.Create(&value).Error
		}
	} else {
		var emptyMessage Message
		//更新消息
		fmt.Printf("更新消息：%s\n", res.MessageId)
		err = global.GDB.Model(&emptyMessage).Select("status").Where("message_id = ?", res.MessageId).Updates(value).Error
	}
	return err
}

func FindMessage(value *Message) (Message, error) {
	var result Message
	err := global.GDB.Where("message_id = ?", value.MessageId).First(&result).Error
	return result, err
}

func FindMessages(whereFrom *Message, whereContact *Message, page int8) (*MessagePageInfo, error) {
	var limit int8
	limit = 50
	//查条数
	var count int8
	global.GDB.Model(&Message{}).Where(whereFrom).Or(whereContact).Count(&count)
	var totalPage float64
	totalPage = float64(count) / float64(limit)
	//偏移量
	offset := count - (page * limit)
	//如果小于0就那之前剩的条数
	if offset < 0 {
		limit = count - ((page - 1) * limit)
	}
	var result []*Message
	err := global.GDB.Where(whereFrom).Or(whereContact).Order("send_time asc").Offset(offset).Limit(50).Find(&result).Error
	messagePageInfo := &MessagePageInfo{
		Page:      page,
		Count:     count,
		PageSize:  limit,
		Data:      result,
		TotalPage: int8(math.Ceil(totalPage)),
	}
	return messagePageInfo, err
}

func GetGroupsMessage(whereFrom *Message) (contactsMessage *ContactsMessage, err error) {
	result := &Message{}
	contactsMessage = &ContactsMessage{}
	err = global.GDB.Where(whereFrom).Order("send_time desc").First(&result).Scan(&contactsMessage).Error
	return
}

type MessagePageInfo struct {
	Page      int8       `json:"page"`
	Count     int8       `json:"count"`
	PageSize  int8       `json:"page_size"`
	Data      []*Message `json:"data"`
	TotalPage int8       `json:"total_page"`
}

func FindGroupMessages(whereFrom *Message, page int8) (*MessagePageInfo, error) {
	var limit int8
	limit = 50
	//查条数
	var count int8
	global.GDB.Model(&Message{}).Where(whereFrom).Count(&count)
	var totalPage float64
	totalPage = float64(count) / float64(limit)
	//偏移量
	offset := count - (page * limit)
	//如果小于0就那之前剩的条数
	if offset < 0 {
		limit = count - ((page - 1) * limit)
	}
	var result []*Message
	err := global.GDB.Where(whereFrom).Order("send_time asc").Offset(offset).Limit(limit).Find(&result).Error
	messagePageInfo := &MessagePageInfo{
		Page:      page,
		Count:     count,
		PageSize:  limit,
		Data:      result,
		TotalPage: int8(math.Ceil(totalPage)),
	}
	return messagePageInfo, err
}
