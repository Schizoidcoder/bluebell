package websocket

import (
	"bluebell/models"
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

func Connect(conn *models.Client) {
	log.Printf("new connection from %s", conn.ID)
	Manager.Clients[conn.ID] = conn
	replyMsg := &models.ReplyMsg{
		From:    "server",
		Code:    200,
		Content: "已连接到服务器",
	}
	msg, _ := json.Marshal(replyMsg)
	_ = conn.Socket.WriteMessage(websocket.TextMessage, msg) //回复给客户端
}

func Disconnect(conn *models.Client) {
	log.Printf("disconnect from %s", conn.ID)
	if _, ok := Manager.Clients[conn.ID]; ok {
		replyMsg := &models.ReplyMsg{
			From:    "server",
			Code:    200,
			Content: "连接已断开",
		}
		msg, _ := json.Marshal(replyMsg)
		_ = conn.Socket.WriteMessage(websocket.TextMessage, msg)
		close(conn.Send)
		delete(Manager.Clients, conn.ID)
	}
}

func CheckIfConnected(Id string) (*models.Client, bool) {
	if conn, ok := Manager.Clients[Id]; ok {
		return conn, true
	}
	return nil, false
}
