package websocket

import (
	"encoding/json"
	"gin-derived/api/models"
	"gin-derived/global"
	"time"
)

type Manager struct {
	//整个客户端
	Clients map[string]*Client
	//系统消息管道
	Broadcast chan []byte
	//在线管道
	Online chan *Client
	//下线管道
	Outline chan *Client
	//普通消息管道
	CommonChan chan []byte
}

var ClientManager = Manager{
	Clients:    make(map[string]*Client),
	Broadcast:  make(chan []byte),
	Online:     make(chan *Client),
	Outline:    make(chan *Client),
	CommonChan: make(chan []byte),
}

func (manager *Manager) Start() {
	for {
		global.GLOG.Infof("开始通信：%s", time.Now().Format("2006-01-02 15:04:05"))
		select {
		case user := <-ClientManager.Online:
			global.GLOG.Infof("用户：%v 加入", user.ID)
			//将新用户放进客户端里
			ClientManager.Clients[user.ID] = user
		case user := <-ClientManager.Outline:
			global.GLOG.Infof("用户：%v 离开", user.ID)
			if _, ok := ClientManager.Clients[user.ID]; ok {
				//关闭发送管道
				close(user.Send)
				//删除客户端
				delete(ClientManager.Clients, user.ID)
			}
		case message := <-ClientManager.Broadcast:
			//对所有人发送的消息管道
			clientMessage := &ClientMessage{}
			err := json.Unmarshal(message, &clientMessage)
			if err != nil {
				global.GLOG.Errorf("发送的信息为空,%v", err.Error())
				break
			}
			HandlerClientMessage(clientMessage, message)
		case message := <-ClientManager.CommonChan:
			//针对某用户发送消息
			Chat := &models.MessageChat{}
			err := json.Unmarshal(message, &Chat)
			if err != nil {
				global.GLOG.Errorf("发送的信息为空,%v", err.Error())
				break
			}
			//查询是否群聊信息
			group := &models.Group{Uuid: Chat.ToContactId}
			info, err := group.FindGroup()
			if err != nil {
				//如果没有就是单聊
				if user, ok := ClientManager.Clients[Chat.ToContactId]; ok {
					//发送时需改变接收用户
					Chat.ToContactId = Chat.FromUser.Id
					NewMessage, _ := json.Marshal(Chat)
					select {
					case user.Send <- NewMessage:
					default:
						close(user.Send)
						delete(ClientManager.Clients, user.ID)
					}
				}
				break
			}
			member, _ := info.FindGroupData()
			for _, v := range member {
				//需要排除是自己发的
				if Chat.FromUser.Id != v.Uuid {
					if user, ok := ClientManager.Clients[v.Uuid]; ok {
						NewMessage, _ := json.Marshal(Chat)
						select {
						case user.Send <- NewMessage:
						default:
							close(user.Send)
							delete(ClientManager.Clients, user.ID)
						}
					}
				}
			}
		}
	}
}

func HandlerClientMessage(message *ClientMessage, msg []byte) {
	switch message.Type {
	case "create_group":
		CreateGroup := &models.CreateGroupMessageResponse{}
		err := json.Unmarshal(msg, &CreateGroup)
		if err != nil {
			global.GLOG.Errorf("发送的信息为空,%v", err.Error())
			break
		}
		//发送给群成员
		if user, ok := ClientManager.Clients[CreateGroup.ToContactId]; ok {
			select {
			case user.Send <- msg:
			default:
				close(user.Send)
				delete(ClientManager.Clients, user.ID)
			}
		}
	case "add_member":
		CreateGroup := &models.CreateGroupMessageResponse{}
		err := json.Unmarshal(msg, &CreateGroup)
		if err != nil {
			global.GLOG.Errorf("发送的信息为空,%v", err.Error())
			break
		}
		//发送给群成员
		if user, ok := ClientManager.Clients[CreateGroup.ToContactId]; ok {
			select {
			case user.Send <- msg:
			default:
				close(user.Send)
				delete(ClientManager.Clients, user.ID)
			}
		}
	}
}
