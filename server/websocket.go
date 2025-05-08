package server

import (
	"log"
	"net/http"
	"time"

	"github.com/Talal52/go-chat/chat/models"
	"github.com/gorilla/websocket"
	"github.com/Talal52/go-chat/chat/service"
	wsmodel "github.com/Talal52/go-chat/chat/models/websocket" // use alias to avoid name conflict

)
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins (you can restrict this in production)
	},
}

func NewWebSocketServer(service *service.ChatService) *wsmodel.WebSocketServer{
    return &wsmodel.WebSocketServer{
        Clients:   make(map[*websocket.Conn]string),
        Broadcast: make(chan models.Message),
        Service:   service,
    }
}

func (server *wsmodel.WebSocketServer) HandleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading connection:", err)
		return
	}
	defer conn.Close()

	var username string
	err = conn.ReadJSON(&username)
	if err != nil {
		log.Println("Error reading username:", err)
		return
	}

	server.Mutex.Lock()
	server.Clients[conn] = username
	server.Mutex.Unlock()

	log.Printf("User %s connected", username)

	for {
		var msg models.Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("Error reading message from %s: %v", username, err)
			break
		}

		msg.Sender = username
		msg.CreatedAt = time.Now()

		if err := server.Service.SaveMessage(msg); err != nil {
			log.Printf("Error saving message: %v", err)
			continue
		}

		server.Broadcast <- msg
	}

	server.Mutex.Lock()
	delete(server.Clients, conn)
	server.Mutex.Unlock()

	log.Printf("User %s disconnected", username)
}


func (server *wsmodel.WebSocketServer) HandleMessages() {
	for {
		msg := <-server.Broadcast

		server.Mutex.Lock()
		for conn := range server.Clients {
			err := conn.WriteJSON(msg)
			if err != nil {
				log.Printf("Error sending message to %s: %v", server.Clients[conn], err)
				conn.Close()
				delete(server.Clients, conn)
			}
		}
		server.Mutex.Unlock()
	}
}
