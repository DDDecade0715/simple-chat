package websocket

import (
	"encoding/json"
	"gin-derived/api/models"
	"gin-derived/api/services/imService"
	"gin-derived/global"
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
	"strings"
	"time"
)

type Client struct {
	ID     string
	Socket *websocket.Conn
	Send   chan []byte
}

type ClientMessage struct {
	Type string `json:"type"`
}

//接收消息
func (c *Client) Read() {
	defer func() {
		ClientManager.Outline <- c
		c.Socket.Close()
	}()

	for {
		c.Socket.PongHandler()
		_, message, err := c.Socket.ReadMessage()
		if err != nil {
			ClientManager.Outline <- c
			c.Socket.Close()
			break
		}
		Message := &ClientMessage{}
		err = json.Unmarshal(message, &Message)
		if err != nil {
			global.GLOG.Errorf("读取到客户端的信息为空,%v", err.Error())
			return
		}
		if err = Message.ReadMessage(message); err != nil {
			global.GLOG.Errorf("读取到客户端的错误信息,%v", err.Error())
			return
		}
	}
}

//ReadMessage 读取收到的消息
func (message *ClientMessage) ReadMessage(msg []byte) error {
	switch message.Type {
	case "create_group":
		//创建群聊
		global.GLOG.Infof("读取到客户端的信息,创建群聊消息,%v", string(msg))
		//向群成员发消息
		CreateGroupMessage := &models.CreateGroupMessage{}
		err := json.Unmarshal(msg, &CreateGroupMessage)
		if err != nil {
			global.GLOG.Errorf("读取到客户端的信息,创建群聊消息为空,%v", err.Error())
			return err
		}
		//查询群信息
		group := &models.Group{
			Uuid: CreateGroupMessage.Id,
		}
		info, err := group.FindGroup()
		if err != nil {
			global.GLOG.Errorf("获取群信息失败,%v", err.Error())
			return err
		}
		members, err := info.FindGroupData()
		if err != nil {
			global.GLOG.Errorf("获取群信息失败,%v", err.Error())
			return err
		}
		var groupUsername string
		var membersUsername string
		for _, v := range members {
			if v.UserId == info.UserId {
				//群主昵称
				groupUsername = v.Username
			} else {
				//成员昵称
				membersUsername += v.Username + "、"
			}
		}
		//保存发送的通知
		num := strings.LastIndex(membersUsername, "、")
		membersUsername = membersUsername[0:num]
		content := "群主 " + groupUsername + " 邀请了 " + membersUsername + " 进入群聊"
		imService.SaveMessage(&models.MessageChat{
			Content: content,
			FromUser: &models.FromUser{
				Avatar:      "",
				DisplayName: "",
				Id:          "",
			},
			SendTime:    time.Now().UnixNano() / 1e6,
			Status:      "succeed",
			ToContactId: group.Uuid,
			Type:        "event",
			Id:          uuid.NewV4().String(),
		})
		for _, v := range members {
			//重写发回前端信息
			Response := &models.CreateGroupMessageResponse{
				Type:        "create_group",
				ToContactId: v.Uuid,
				GroupInfo: &models.GroupInfo{
					Id:     info.ID,
					Uuid:   info.Uuid,
					Name:   info.Name,
					Avatar: info.Avatar,
					Message: &models.ContactsMessage{
						Type:     "event",
						Content:  content,
						SendTime: time.Now().Format("2006-01-02 15:04:05"),
						Status:   "succeed",
					},
					Members: members,
				},
			}
			response, err := json.Marshal(Response)
			if err != nil {
				global.GLOG.Errorf("返回信息失败,%v", err.Error())
				return err
			}
			//系统消息
			ClientManager.Broadcast <- response
		}
	case "add_member":
		//新增群成员
		global.GLOG.Infof("读取到客户端的信息,新增群成员消息,%v", string(msg))
		//向群成员发消息
		CreateGroupMessage := &models.CreateGroupMessage{}
		err := json.Unmarshal(msg, &CreateGroupMessage)
		if err != nil {
			global.GLOG.Errorf("读取到客户端的信息,新增群成员消息为空,%v", err.Error())
			return err
		}
		//查询群信息
		group := &models.Group{
			Uuid: CreateGroupMessage.Id,
		}
		info, err := group.FindGroup()
		if err != nil {
			global.GLOG.Errorf("获取群信息失败,%v", err.Error())
			return err
		}
		members, err := info.FindGroupData()
		if err != nil {
			global.GLOG.Errorf("获取群信息失败,%v", err.Error())
			return err
		}
		var groupUsername string
		var membersUsername string
		//被邀请的成员
		for _, v := range members {
			if v.UserId == info.UserId {
				//群主昵称
				groupUsername = v.Username
			}
			for _, value := range CreateGroupMessage.Members {
				if value == v.Uuid {
					//成员昵称
					membersUsername += v.Username + "、"
				}
			}
		}
		//保存发送的通知
		num := strings.LastIndex(membersUsername, "、")
		membersUsername = membersUsername[0:num]
		content := "群主 " + groupUsername + " 邀请了 " + membersUsername + " 进入群聊"
		messageChat := &models.MessageChat{
			Content: content,
			FromUser: &models.FromUser{
				Avatar:      group.Avatar,
				DisplayName: group.Name,
				Id:          group.Uuid,
			},
			SendTime:    time.Now().UnixNano() / 1e6,
			Status:      "succeed",
			ToContactId: group.Uuid,
			Type:        "event",
			Id:          uuid.NewV4().String(),
		}
		imService.SaveMessage(messageChat)

		for _, v := range members {
			//重写发回前端信息
			Response := &models.CreateGroupMessageResponse{
				Type:        "add_member",
				ToContactId: v.Uuid,
				GroupInfo: &models.GroupInfo{
					Id:     info.ID,
					Uuid:   info.Uuid,
					Name:   info.Name,
					Avatar: info.Avatar,
					Message: &models.ContactsMessage{
						Type:     "event",
						Content:  content,
						SendTime: time.Now().Format("2006-01-02 15:04:05"),
						Status:   "succeed",
					},
					Members: members,
				},
			}
			response, err := json.Marshal(Response)
			if err != nil {
				global.GLOG.Errorf("返回信息失败,%v", err.Error())
				return err
			}
			//系统消息
			ClientManager.Broadcast <- response

			messageChat.ToContactId = v.Uuid
			messageC, err := json.Marshal(messageChat)
			if err != nil {
				global.GLOG.Errorf("返回信息失败,%v", err.Error())
				return err
			}
			ClientManager.CommonChan <- messageC
		}
	case "edit_userinfo":

	case "text":
		global.GLOG.Infof("读取到客户端的信息,用户聊天消息,%v", string(msg))
		Chat := &models.MessageChat{}
		err := json.Unmarshal(msg, &Chat)
		if err != nil {
			global.GLOG.Errorf("读取到客户端的信息,用户聊天消息为空,%v", err.Error())
			return err
		}
		imService.SaveMessage(Chat)
		//用户聊天消息管道
		ClientManager.CommonChan <- msg
	case "image":
		global.GLOG.Infof("读取到客户端的信息,用户聊天图片消息,%v", string(msg))
		Chat := &models.MessageChat{}
		err := json.Unmarshal(msg, &Chat)
		if err != nil {
			global.GLOG.Errorf("读取到客户端的信息,用户聊天图片消息为空,%v", err.Error())
			return err
		}
		imService.SaveMessage(Chat)
		//用户聊天消息管道
		ClientManager.CommonChan <- msg
	case "file":
		global.GLOG.Infof("读取到客户端的信息,用户文件消息,%v", string(msg))
		Chat := &models.MessageChat{}
		err := json.Unmarshal(msg, &Chat)
		if err != nil {
			global.GLOG.Errorf("读取到客户端的信息,用户文件消息为空,%v", err.Error())
			return err
		}
		imService.SaveMessage(Chat)
		//用户聊天消息管道
		ClientManager.CommonChan <- msg
	default:
		global.GLOG.Infof("读取到客户端的信息,没有规定类型,%v", string(msg))
	}
	return nil
}

//发送消息
func (c *Client) Write() {
	defer func() {
		c.Socket.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			//这里是拿消息失败
			if !ok {
				c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			global.GLOG.Infof("发送到%s客户端的信息:%s\n", c.ID, string(message))
			//发送文字消息
			c.Socket.WriteMessage(websocket.TextMessage, message)
		}
	}
}
