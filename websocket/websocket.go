package websocket

import "bluebell/models"

var Manager = &models.ClientManager{
	Clients:    make(map[string]*models.Client),
	Broadcast:  make(chan *models.Broadcast),
	Reply:      make(chan *models.Client),
	Register:   make(chan *models.Client),
	Unregister: make(chan *models.Client),
}

func Init() {
	for {
		select {
		case conn := <-Manager.Register:
			Connect(conn)
		case conn := <-Manager.Unregister:
			Disconnect(conn)

		}
	}
}
