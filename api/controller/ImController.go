package controller

import (
	"gin-derived/api/services/imService"
	"github.com/gin-gonic/gin"
)

func ImLogin(c *gin.Context) {
	imService.LoginFromIm(c)
}

func GetUserContacts(c *gin.Context) {
	imService.GetUserContacts(c)
}

func GetMessages(c *gin.Context) {
	imService.GetMessages(c)
}

func UploadAvatar(c *gin.Context) {
	imService.UploadAvatar(c)
}

func UpdateUserinfo(c *gin.Context) {
	imService.UpdateUserinfo(c)
}

func UploadChatImage(c *gin.Context) {
	imService.UploadChatImage(c)
}

func GetMessageInfoById(c *gin.Context) {
	imService.GetMessageInfoById(c)
}
