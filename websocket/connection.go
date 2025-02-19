package websocket

import (
	"bluebell/dao/mongo"
	"bluebell/models"
	"encoding/json"
	"log"
	"strconv"

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
	doc := make(map[string]interface{})
	intID, err2 := strconv.ParseInt(conn.ID, 10, 64)
	if err2 != nil {
		log.Fatalln(err2)
		return
	}
	doc["AuthorID"] = intID
	results, err := mongo.FindManyByOneCon("message", doc)
	if err != nil {
		log.Println(err)
	} else {
		marshal, err := json.Marshal(results)
		if err != nil {
			log.Println(err)
		} else {
			//fmt.Println(results)
			//fmt.Println(string(marshal))
			_ = conn.Socket.WriteMessage(websocket.TextMessage, marshal)
		}
	}

	doc2 := make(map[string]interface{})
	doc2["Recipient"] = conn.ID
	results2, err2 := mongo.FindManyByOneCon("message", doc2)
	if err2 != nil {
		log.Println(err2)
	} else {
		marshal2, err := json.Marshal(results2)
		if err != nil {
			log.Println(err)
		} else {
			//fmt.Println(results)
			//fmt.Println(string(marshal))
			_ = conn.Socket.WriteMessage(websocket.TextMessage, marshal2)
		}
	}

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
