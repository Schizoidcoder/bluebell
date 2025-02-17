package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
	mywebsocket "bluebell/websocket"
	"strconv"

	"github.com/goccy/go-json"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

func ConnectWebSocket(id string, conn *websocket.Conn) {
	// 创建一个用户实例
	client := &models.Client{
		ID:     id,
		Socket: conn,
		Send:   make(chan []byte),
	}
	// 用户注册到用户管理上
	mywebsocket.Manager.Register <- client
	go Read(client)
	go Write(client)
}

func Read(c *models.Client) {
	defer func() { //避免忘记关闭，所以要加上close
		mywebsocket.Manager.Unregister <- c
		_ = c.Socket.Close()
	}()
	for {
		//c.Socket.PongHandler() //心跳机制，确保连接正常 //socket会自动调用
		SendMessage := &models.SendMessage{}
		err := c.Socket.ReadJSON(SendMessage)
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				// 正常关闭，不打印为数据格式错误
				_ = c.Socket.Close()
				break
			}
			zap.L().Error("数据格式不正确", zap.Error(err))
			mywebsocket.Manager.Unregister <- c
			_ = c.Socket.Close()
			break
		}
		if SendMessage.Recipient == "" {
			replyMsg := models.ReplyMsg{
				From:    "server",
				Code:    0,
				Content: "用户不存在",
			}
			msg, _ := json.Marshal(replyMsg)
			_ = c.Socket.WriteMessage(websocket.TextMessage, msg)
			continue
		}
		Numid, err := strconv.ParseInt(SendMessage.Recipient, 10, 64)
		if err != nil {
			replyMsg := models.ReplyMsg{
				From:    "server",
				Code:    0,
				Content: "用户不存在",
			}
			msg, _ := json.Marshal(replyMsg)
			_ = c.Socket.WriteMessage(websocket.TextMessage, msg)
			continue
		}
		_, err = mysql.CheckUserExistById(Numid)
		if err != nil {
			replyMsg := models.ReplyMsg{
				From:    "server",
				Code:    0,
				Content: "用户不存在",
			}
			msg, _ := json.Marshal(replyMsg)
			_ = c.Socket.WriteMessage(websocket.TextMessage, msg)
			continue
		}
		AimClient, flag := mywebsocket.CheckIfConnected(SendMessage.Recipient)
		if !flag {
			replyMsg := models.ReplyMsg{
				From:    "server",
				Code:    0,
				Content: "用户已掉线",
			}
			msg, _ := json.Marshal(replyMsg)
			_ = c.Socket.WriteMessage(websocket.TextMessage, msg)
			//Todo：用户下线后的逻辑
			continue
		}
		replyMsg := models.ReplyMsg{
			From:    c.ID,
			Code:    0,
			Content: SendMessage.Content,
		}
		var marshalreply []byte
		marshalreply, err = json.Marshal(replyMsg)
		AimClient.Send <- marshalreply
	}
}

func Write(c *models.Client) {
	defer func() {
		_ = c.Socket.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				_ = c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			zap.L().Debug(c.ID + " 接受消息:" + string(message))
			_ = c.Socket.WriteMessage(websocket.TextMessage, message)
		}
	}

}
