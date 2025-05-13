package server

import (
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/Talal52/go-chat/chat/models"
	"github.com/Talal52/go-chat/chat/service"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var upgrader = websocket.Upgrader{}

func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	tokenString, err := extractToken(r.Header.Get("Authorization"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	sender, err := parseToken(tokenString)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}
		log.Printf("Received message from %s: %s", sender, msg)
	}
}

func extractToken(authHeader string) (string, error) {
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return "", http.ErrNoCookie
	}
	return strings.TrimPrefix(authHeader, "Bearer "), nil
}

func parseToken(tokenString string) (string, error) {
	claims := &jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil || !token.Valid {
		return "", http.ErrNoCookie
	}

	sender, ok := (*claims)["username"].(string)
	if !ok {
		return "", http.ErrNoCookie
	}
	return sender, nil
}

type WebSocketServer struct {
	Clients   map[*websocket.Conn]string
	Broadcast chan models.Message
	Service   *service.ChatService
	Mutex     sync.Mutex
}

func NewWebSocketServer(service *service.ChatService) *WebSocketServer {
	return &WebSocketServer{
		Clients:   make(map[*websocket.Conn]string),
		Broadcast: make(chan models.Message),
		Service:   service,
	}
}

func (server *WebSocketServer) HandleConnections(w http.ResponseWriter, r *http.Request) {
	// Extract and validate JWT token from Authorization header
	tokenString, err := extractToken(r.Header.Get("Authorization"))
	if err != nil {
		http.Error(w, "Unauthorized: Missing or invalid token", http.StatusUnauthorized)
		return
	}

	sender, err := parseToken(tokenString)
	if err != nil {
		http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
		return
	}

	// Upgrade the HTTP connection to a WebSocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading connection:", err)
		return
	}
	defer conn.Close()

	server.Mutex.Lock()
	server.Clients[conn] = sender
	server.Mutex.Unlock()

	log.Printf("User %s connected", sender)

	for {
		var msg models.Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("Error reading message from %s: %v", sender, err)
			break
		}

		msg.Sender = sender
		msg.CreatedAt = time.Now()

		if err := server.Service.SaveMessage(msg); err != nil {
			log.Printf("Error saving message: %v", err)
			continue
		}

		if msg.GroupID != nil {
			// Send to all connected users in the group
			server.BroadcastToGroup(*msg.GroupID, msg)
		} else {
			server.SendMessageToRecipient(msg)
		}
	}

	server.Mutex.Lock()
	delete(server.Clients, conn)
	server.Mutex.Unlock()

	log.Printf("User %s disconnected", sender)
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

func (server *WebSocketServer) SendMessageToRecipient(msg models.Message) {
	server.Mutex.Lock()
	defer server.Mutex.Unlock()

	for conn, username := range server.Clients {
		if username == msg.Receiver {
			err := conn.WriteJSON(msg)
			if err != nil {
				log.Printf("Error sending message to %s: %v", username, err)
				conn.Close()
				delete(server.Clients, conn)
			}
			return
		}
	}

	log.Printf("Recipient %s not connected", msg.Receiver)
}
func (server *WebSocketServer) BroadcastToGroup(groupID primitive.ObjectID, msg models.Message) {
	server.Mutex.Lock()
	defer server.Mutex.Unlock()

	for conn, username := range server.Clients {
		// For now, we'll broadcast to all for simplicity
		if err := conn.WriteJSON(msg); err != nil {
			log.Printf("Error sending group message to %s: %v", username, err)
			conn.Close()
			delete(server.Clients, conn)
		}
	}
}
