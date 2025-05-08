package server

import (
    "log"
    "net/http"
    "sync"
    "time"

    "github.com/Talal52/go-chat/chat/models"
    // wsmodel "github.com/Talal52/go-chat/chat/models/websocket" // Alias to avoid name conflict
    "github.com/Talal52/go-chat/chat/service"
    "github.com/gorilla/websocket"
)

type WebSocketServer struct {
    Clients   map[*websocket.Conn]string // Map of WebSocket connections to usernames
    Broadcast chan models.Message        // Channel for broadcasting messages
    Service   *service.ChatService       // Chat service for saving messages
    Mutex     sync.Mutex                 // Mutex to protect the Clients map
}

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true // Allow all origins (you can restrict this in production)
    },
}

// NewWebSocketServer initializes a new WebSocket server
func NewWebSocketServer(service *service.ChatService) *WebSocketServer {
    return &WebSocketServer{
        Clients:   make(map[*websocket.Conn]string),
        Broadcast: make(chan models.Message),
        Service:   service,
    }
}

// HandleConnections handles incoming WebSocket connections
func (server *WebSocketServer) HandleConnections(w http.ResponseWriter, r *http.Request) {
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

// HandleMessages handles broadcasting messages to all connected clients
func (server *WebSocketServer) HandleMessages() {
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