package server

import (
	"log"
	"net/http"
	"time"

	"github.com/Talal52/go-chat/chat/models"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins (you can restrict this in production)
	},
}

func NewWebSocketServer(service *models.ChatService) *websocket.WebSocketServer {
	return &websocket.WebSocketServer{
		Clients:   make(map[*websocket.Conn]string),
		Broadcast: make(chan models.Message),
		Service:   service,
	}
}

func (ws *websocket.WebSocketServer) HandleConnections(w http.ResponseWriter, r *http.Request) {
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

	ws.Mutex.Lock()
	ws.Clients[conn] = username
	ws.Mutex.Unlock()

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

		if err := ws.Service.SaveMessage(msg); err != nil {
			log.Printf("Error saving message: %v", err)
			continue
		}

		ws.Broadcast <- msg
	}

	ws.Mutex.Lock()
	delete(ws.Clients, conn)
	ws.Mutex.Unlock()

	log.Printf("User %s disconnected", username)
}

func (ws *websocket.WebSocketServer) HandleMessages() {
	for {
		msg := <-ws.Broadcast

		ws.Mutex.Lock()
		for conn := range ws.Clients {
			err := conn.WriteJSON(msg)
			if err != nil {
				log.Printf("Error sending message to %s: %v", ws.Clients[conn], err)
				conn.Close()
				delete(ws.Clients, conn)
			}
		}
		ws.Mutex.Unlock()
	}
}
