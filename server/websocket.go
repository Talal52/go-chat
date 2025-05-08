package server

import (
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/Talal52/go-chat/chat/models"
	"github.com/golang-jwt/jwt"
	"github.com/Talal52/go-chat/chat/service"
	"github.com/gorilla/websocket"
)

type WebSocketServer struct {
	Clients   map[*websocket.Conn]string
	Broadcast chan models.Message
	Service   *service.ChatService
	Mutex     sync.Mutex
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NewWebSocketServer(service *service.ChatService) *WebSocketServer {
	return &WebSocketServer{
		Clients:   make(map[*websocket.Conn]string),
		Broadcast: make(chan models.Message),
		Service:   service,
	}
}

func (server *WebSocketServer) HandleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading connection:", err)
		return
	}
	defer conn.Close()

	var tokenString string
	if err := conn.ReadJSON(&tokenString); err != nil {
		log.Println("Error reading token:", err)
		return
	}

	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil || !token.Valid {
		log.Println("Invalid token:", err)
		return
	}

	username, ok := claims["username"].(string)
	if !ok {
		log.Println("Invalid token claims: username not found")
		return
	}

	server.Mutex.Lock()
	server.Clients[conn] = username
	server.Mutex.Unlock()

	log.Printf("User %s connected", username)

	for {
		var msg models.Message
		if err := conn.ReadJSON(&msg); err != nil {
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

func (server *WebSocketServer) HandleMessages() {
	for msg := range server.Broadcast {
		server.Mutex.Lock()
		for conn, username := range server.Clients {
			if err := conn.WriteJSON(msg); err != nil {
				log.Printf("Error sending message to %s: %v", username, err)
				conn.Close()
				delete(server.Clients, conn)
			}
		}
		server.Mutex.Unlock()
	}
}
