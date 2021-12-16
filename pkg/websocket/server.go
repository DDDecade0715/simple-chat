package websocket

import (
	"gin-derived/api/models"
	jwtPkg "gin-derived/pkg/jwt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

//WsHandler socket 连接 中间件 作用:升级协议,用户验证,自定义信息等
func WsHandler(c *gin.Context) {
	ws := &websocket.Upgrader{
		//设置允许跨域
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		//设置请求协议
		Subprotocols: []string{c.GetHeader("Sec-WebSocket-Protocol")},
	}
	conn, err := ws.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		http.NotFound(c.Writer, c.Request)
		return
	}
	//处理用户信息并把用户注册到客户端
	token := c.Query("token")
	claims, _ := jwtPkg.ParseToken(token)
	user, err := models.FindUser(&models.User{
		ID: uint(claims.UserID),
	})
	if err != nil {
		http.NotFound(c.Writer, c.Request)
		return
	}
	//可以添加用户信息验证
	client := &Client{
		ID:     user.Uuid,
		Socket: conn,
		Send:   make(chan []byte),
	}
	ClientManager.Online <- client
	go client.Read()
	go client.Write()
}
