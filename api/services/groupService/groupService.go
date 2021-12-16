package groupService

import (
	"gin-derived/api/models"
	"gin-derived/pkg/app/response"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

type CreateGroupForm struct {
	ContactIds []string `json:"contact_ids" binding:"required"`
}

//CreateGroup 创建群聊
func CreateGroup(c *gin.Context) {
	var form CreateGroupForm
	if err := c.Bind(&form); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	userId, _ := c.Get("UserID")
	value := &models.User{
		ID: uint(userId.(int64)),
	}
	user, err := models.FindUser(value)
	if err != nil {
		response.FailWithMessage("fail", c)
		return
	}
	group := &models.Group{
		UserId: user.ID,
		Name:   user.Username + "建的群",
		Uuid:   uuid.NewV4().String(),
		Avatar: "",
	}
	form.ContactIds = append(form.ContactIds, user.Uuid)
	err = group.CreateGroup(form.ContactIds)
	if err != nil {
		response.FailWithMessage("创建群组失败", c)
		return
	}
	response.OkWithData("创建群组成功", group, c)
	return
}

//GetGroupContacts 获取所有群聊
func GetGroupContacts(c *gin.Context) {
	userId, _ := c.Get("UserID")
	groupMember := &models.GroupMember{
		UserId: uint(userId.(int64)),
	}
	groupIds, err := groupMember.FindGroups()
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	var ids []uint
	for _, v := range groupIds {
		ids = append(ids, v.GroupId)
	}
	group := &models.Group{}
	res, err := group.FindGroupInfo(ids)
	if err != nil {
		response.FailWithMessage("fail", c)
		return
	}
	for _, v := range res {
		//获取群聊天记录
		whereFrom := &models.Message{
			ToContactId: v.Uuid,
		}
		m, _ := models.GetGroupsMessage(whereFrom)
		v.Message = m
		//获取群成员
		group.Uuid = v.Uuid
		data, _ := group.FindGroupData()
		v.Members = data
	}
	response.OkWithData("success", res, c)
	return
}

type AddGroupMemberRequest struct {
	GroupId    string   `json:"group_id" binding:"required"`
	ContactIds []string `json:"contact_ids" binding:"required"`
}

//AddGroupMembers 添加成员入群聊
func AddGroupMembers(c *gin.Context) {
	var form AddGroupMemberRequest
	if err := c.Bind(&form); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	group := &models.Group{
		Uuid: form.GroupId,
	}
	groupInfo, err := group.FindGroup()
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	for _, v := range form.ContactIds {
		groupMember := &models.GroupMember{
			GroupId: groupInfo.ID,
		}
		_, err = groupMember.AddGroupMember(v)
		if err != nil {
			response.FailWithMessage(err.Error(), c)
			return
		}
	}
	response.OkWithData("success", groupInfo, c)
	return
}

type GetGroupMembersRequest struct {
	GroupId string `json:"group_id" binding:"required"`
}

//GetGroupMembers 获取群成员
func GetGroupMembers(c *gin.Context) {
	var form GetGroupMembersRequest
	if err := c.Bind(&form); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	group := &models.Group{
		Uuid: form.GroupId,
	}
	res, err := group.FindGroupData()
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData("success", res, c)
	return
}

type UpdateGroup struct {
	GroupId string `json:"group_id" binding:"required"`
	Avatar  string `json:"avatar" binding:"required"`
	Name    string `json:"name" binding:"required"`
}

func UploadGroup(c *gin.Context) {
	var form UpdateGroup
	if err := c.Bind(&form); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	group := &models.Group{
		Uuid: form.GroupId,
	}
	update := &models.Group{
		Name:   form.Name,
		Avatar: form.Avatar,
	}
	res, err := group.UpdateGroup(update)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData("success", res, c)
	return
}
