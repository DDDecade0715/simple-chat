package controller

import (
	"gin-derived/api/services/groupService"
	"github.com/gin-gonic/gin"
)

func CreateGroup(c *gin.Context) {
	groupService.CreateGroup(c)
}

func GetGroupContacts(c *gin.Context) {
	groupService.GetGroupContacts(c)
}

func AddGroupMembers(c *gin.Context) {
	groupService.AddGroupMembers(c)
}

func GetGroupMembers(c *gin.Context) {
	groupService.GetGroupMembers(c)
}

func UploadGroup(c *gin.Context) {
	groupService.UploadGroup(c)
}
