package imService

import (
	"encoding/json"
	"gin-derived/api/models"
	"gin-derived/global"
	"gin-derived/pkg/app/response"
	"gin-derived/pkg/jwt"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"io"
	"os"
	"time"
)

type LoginStruct struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

type LoginResponse struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Token    string `json:"token"`
	Avatar   string `json:"avatar"`
}

func LoginFromIm(c *gin.Context) {
	var form LoginStruct
	if err := c.Bind(&form); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	res, err := models.FindUser(&models.User{
		Username: form.Username,
	})
	if err != nil {
		//找不到直接创建
		res, err = models.CreateUser(&models.User{
			Username: form.Username,
			Password: form.Password,
			Uuid:     uuid.NewV4().String(),
		})
		if err != nil {
			response.FailWithMessage("登录失败", c)
			return
		}
	}
	if res.Password != form.Password {
		response.FailWithMessage("密码不正确", c)
		return
	}
	id := int64(res.ID)
	token, _ := jwt.GenerateToken(id, res.Username)
	resp := &LoginResponse{
		Id:       res.Uuid,
		Username: res.Username,
		Token:    token,
		Avatar:   res.Avatar,
	}
	response.OkWithData("登录成功", resp, c)
	return
}

//GetUserContacts 获取联系人列表
func GetUserContacts(c *gin.Context) {
	userId, _ := c.Get("UserID")
	value := &models.User{
		ID: uint(userId.(int64)),
	}
	//查询用户信息
	user, err := models.FindUser(value)
	if err != nil {
		response.FailWithMessage("fail", c)
		return
	}
	//获取好友列表
	res, _ := models.FindUserContacts(value)
	//查询好友聊天记录
	for _, v := range res {
		whereFrom := &models.Message{
			FromUserId:  user.Uuid,
			ToContactId: v.Uuid,
		}
		whereContact := &models.Message{
			FromUserId:  v.Uuid,
			ToContactId: user.Uuid,
		}
		m, _ := models.GetContactsMessage(whereFrom, whereContact)
		v.Message = m
	}

	response.OkWithData("success", res, c)
	return
}

//SaveMessage 保存聊天记录
func SaveMessage(value *models.MessageChat) {
	var fromUser []byte
	if value.FromUser != (&models.FromUser{}) {
		fromUser, _ = json.Marshal(value.FromUser)
	}
	tm := time.Unix(value.SendTime/1000, 0)

	message := &models.Message{
		MessageId:   value.Id,
		FromUserId:  value.FromUser.Id,
		ToContactId: value.ToContactId,
		Type:        value.Type,
		Content:     value.Content,
		FromUser:    string(fromUser),
		Status:      value.Status,
		SendTime:    tm.Format("2006-01-02 15:04:05"),
	}
	err := models.SaveOrUpdateMessage(message)
	if err != nil {
		global.GLOG.Infof("创建消息失败:%d\n", err)
	}
}

type ChatContact struct {
	ContactId string `json:"contact_id" binding:"required"`
	Type      string `json:"type" binding:"required"`
}

//GetMessages 获取联系人聊天记录
func GetMessages(c *gin.Context) {
	var form ChatContact
	userId, _ := c.Get("UserID")
	if err := c.Bind(&form); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	switch form.Type {
	case "group":
		whereFrom := &models.Message{
			ToContactId: form.ContactId,
		}
		res, _ := models.FindGroupMessages(whereFrom)
		var result []*models.MessageChat
		for _, v := range res {
			//发送人组装
			FromByte := []byte(v.FromUser)
			var From *models.FromUser
			var fileName string
			var fileSize int64
			json.Unmarshal(FromByte, &From)
			//时间转换
			timer, _ := time.ParseInLocation("2006-01-02 15:04:05", v.SendTime, time.Local)

			if v.Type == "image" {
				chatImage, _ := models.FindChatImage(&models.ChatImage{
					MessageId: v.MessageId,
				})
				v.Content = chatImage.Url
			}

			if v.Type == "file" {
				chatImage, _ := models.FindChatImage(&models.ChatImage{
					MessageId: v.MessageId,
				})
				v.Content = chatImage.Url
				fileName = chatImage.FileName
				fileSize = chatImage.FileSize
			}

			MessagesChat := &models.MessageChat{
				Id:          v.MessageId,
				Content:     v.Content,
				FromUser:    From,
				SendTime:    timer.UnixNano() / 1e6,
				Status:      v.Status,
				ToContactId: v.ToContactId,
				Type:        v.Type,
				MessageId:   v.MessageId,
				FileName:    fileName,
				FileSize:    fileSize,
			}
			result = append(result, MessagesChat)
		}
		response.OkWithData("success", result, c)
		return
	case "user":
		value := &models.User{}
		value.ID = uint(userId.(int64))
		user, err := models.FindUser(value)
		if err != nil {
			response.FailWithMessage("fail", c)
			return
		}

		whereFrom := &models.Message{
			FromUserId:  user.Uuid,
			ToContactId: form.ContactId,
		}
		whereContact := &models.Message{
			FromUserId:  form.ContactId,
			ToContactId: user.Uuid,
		}
		res, _ := models.FindMessages(whereFrom, whereContact)
		var result []*models.MessageChat
		for _, v := range res {
			//发送人组装
			FromByte := []byte(v.FromUser)
			var From *models.FromUser
			var fileName string
			var fileSize int64
			json.Unmarshal(FromByte, &From)
			//时间转换
			timer, _ := time.ParseInLocation("2006-01-02 15:04:05", v.SendTime, time.Local)

			if v.Type == "image" {
				chatImage, _ := models.FindChatImage(&models.ChatImage{
					MessageId: v.MessageId,
				})
				v.Content = chatImage.Url
			}

			if v.Type == "file" {
				chatImage, _ := models.FindChatImage(&models.ChatImage{
					MessageId: v.MessageId,
				})
				v.Content = chatImage.Url
				fileName = chatImage.FileName
				fileSize = chatImage.FileSize
			}

			MessagesChat := &models.MessageChat{
				Id:          v.MessageId,
				Content:     v.Content,
				FromUser:    From,
				SendTime:    timer.UnixNano() / 1e6,
				Status:      v.Status,
				ToContactId: v.ToContactId,
				Type:        v.Type,
				MessageId:   v.MessageId,
				FileName:    fileName,
				FileSize:    fileSize,
			}
			result = append(result, MessagesChat)
		}
		response.OkWithData("success", result, c)
		return
	}
}

type FileInfo struct {
	Name      string
	AccessUrl string
}

//UploadAvatar 上传图像
func UploadAvatar(c *gin.Context) {
	_, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	//判断上传地址是否存在
	savePath := global.GCONFIG.App.UploadPath
	_, osErr := os.Stat(savePath)
	if os.IsNotExist(osErr) {
		err = os.MkdirAll(savePath, os.ModePerm)
		if err != nil {
			response.FailWithMessage("failed to create save directory.", c)
			return
		}
	}
	//判断权限
	if os.IsPermission(osErr) {
		response.FailWithMessage("insufficient file permissions.", c)
		return
	}

	fileName := fileHeader.Filename
	dst := savePath + "/" + fileName
	//保存文件
	src, err := fileHeader.Open()
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithData("success.", &FileInfo{
		Name:      fileName,
		AccessUrl: global.GCONFIG.App.UploadUrl + "/" + fileName,
	}, c)
	return
}

type UserInfo struct {
	Account string `json:"account" binding:"required"`
	Avatar  string `json:"avatar" binding:"required"`
}

//UpdateUserinfo 更新用户资料
func UpdateUserinfo(c *gin.Context) {
	var form UserInfo
	userId, _ := c.Get("UserID")
	if err := c.Bind(&form); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	_, err := models.FindUsername(&models.User{
		ID:       uint(userId.(int64)),
		Username: form.Account,
	})
	if err == nil {
		response.FailWithMessage("用户账号存在重复", c)
		return
	}
	err = models.UpdateUserinfo(&models.User{
		ID:       uint(userId.(int64)),
		Username: form.Account,
		Avatar:   form.Avatar,
	})
	if err != nil {
		response.FailWithMessage("fail", c)
		return
	}
	response.OkWithData("success", form, c)
}

//UploadChatImage 上传聊天图片
func UploadChatImage(c *gin.Context) {
	_, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	//判断上传地址是否存在
	savePath := global.GCONFIG.App.UploadPath
	realPath := "/" + time.Now().Format("20060102")
	//按日期分组
	savePath += realPath
	_, osErr := os.Stat(savePath)
	if os.IsNotExist(osErr) {
		err = os.MkdirAll(savePath, os.ModePerm)
		if err != nil {
			response.FailWithMessage("failed to create save directory.", c)
			return
		}
	}
	//判断权限
	if os.IsPermission(osErr) {
		response.FailWithMessage("insufficient file permissions.", c)
		return
	}

	var fileSize int64
	fileName := fileHeader.Filename
	dst := savePath + "/" + fileName
	//保存文件
	src, err := fileHeader.Open()
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	fi, err := os.Stat(dst)
	if err == nil {
		fileSize = fi.Size()
	}

	chatId, ok := c.GetPostForm("chat_id")
	if !ok {
		response.FailWithMessage("fail", c)
		return
	}

	value := &models.ChatImage{
		MessageId: chatId,
		Url:       global.GCONFIG.App.UploadUrl + realPath + "/" + fileName,
		FileSize:  fileSize,
		FileName:  fileName,
	}
	err = models.CreateChatImage(value)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData("success", value, c)
}

type FormMessage struct {
	MessageId string `json:"message_id" binding:"required"`
}

//GetMessageInfoById 通过聊天记录ID获取url
func GetMessageInfoById(c *gin.Context) {
	var form FormMessage
	if err := c.Bind(&form); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	chatVideo, _ := models.FindChatImage(&models.ChatImage{
		MessageId: form.MessageId,
	})
	result := &models.ChatImage{
		MessageId: chatVideo.MessageId,
		Url:       chatVideo.Url,
	}
	response.OkWithData("success", result, c)
}
