package websocket

import (
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/Talal52/go-chat/chat/models"
	"github.com/Talal52/go-chat/chat/service"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for development
	},
}

type Client struct {
	ID   string
	Conn *websocket.Conn
	Send chan models.Message
}

type WebSocketServer struct {
	Clients   map[string]*Client // Map userID to client
	Broadcast chan models.Message
	Service   *service.ChatService
	Mutex     sync.Mutex
}

func NewWebSocketServer(service *service.ChatService) *WebSocketServer {
	return &WebSocketServer{
		Clients:   make(map[string]*Client),
		Broadcast: make(chan models.Message),
		Service:   service,
	}
}

func (server *WebSocketServer) HandleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading connection:", err)
		return
	}
	defer conn.Close()

	// Read userID from query
	userID := r.URL.Query().Get("userId")
	if userID == "" {
		log.Println("Missing userId")
		return
	}

	client := &Client{
		ID:   userID,
		Conn: conn,
		Send: make(chan models.Message),
	}
	server.Mutex.Lock()
	server.Clients[userID] = client
	server.Mutex.Unlock()

	log.Println("New WebSocket connection established for user:", userID)

	go server.writePump(client)
	server.readPump(client)
}

func (server *WebSocketServer) readPump(client *Client) {
    defer func() {
        server.Mutex.Lock()
        client.Conn.Close()
        delete(server.Clients, client.ID)
        server.Mutex.Unlock()
        log.Println("WebSocket connection closed for user:", client.ID)
    }()

    for {
        var payload struct {
            Message    string `json:"message"`
            ReceiverId int    `json:"receiver_id"`
        }
        err := client.Conn.ReadJSON(&payload)
        if err != nil {
            log.Println("Read error:", err)
            break
        }

        msg := models.Message{
            SenderID:   client.ID,
            Message:    payload.Message,
            ReceiverID: strconv.Itoa(payload.ReceiverId), // Convert int to string
            Timestamp:  time.Now(),
        }

        if msg.ReceiverID == "" {
            log.Println("Missing receiverId in message")
            continue
        }

        if err := server.Service.SaveMessage(msg); err != nil {
            log.Println("Error saving message:", err)
            continue
        }

        server.Mutex.Lock()
        if receiverClient, ok := server.Clients[msg.ReceiverID]; ok {
            receiverClient.Send <- msg
        } else {
            log.Println("Receiver not connected:", msg.ReceiverID)
        }
        server.Mutex.Unlock()
    }
}

func (server *WebSocketServer) writePump(client *Client) {
	defer func() {
		server.Mutex.Lock()
		client.Conn.Close()
		server.Mutex.Unlock()
	}()

	for {
		msg, ok := <-client.Send
		if !ok {
			client.Conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}
		err := client.Conn.WriteJSON(msg)
		if err != nil {
			log.Println("Error sending message:", err)
			return
		}
	}
}

func (server *WebSocketServer) HandleMessages() {
	for msg := range server.Broadcast {
		server.Mutex.Lock()
		for _, client := range server.Clients {
			select {
			case client.Send <- msg:
			default:
				log.Println("Failed to send to client:", client.ID)
			}
		}
		server.Mutex.Unlock()
	}
}