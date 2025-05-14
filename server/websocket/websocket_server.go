package websocket

import (
    "log"
    "net/http"
    "sync"

    "github.com/Talal52/go-chat/chat/models"
    "github.com/Talal52/go-chat/chat/service"
    "github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

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
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println("Error upgrading connection:", err)
        return
    }
    defer conn.Close()

    server.Mutex.Lock()
    server.Clients[conn] = "anonymous" // Replace with actual user identification if needed
    server.Mutex.Unlock()

    log.Println("New WebSocket connection established")

    for {
        var msg models.Message
        err := conn.ReadJSON(&msg)
        if err != nil {
            log.Println("Error reading message:", err)
            break
        }

        server.Broadcast <- msg
    }

    server.Mutex.Lock()
    delete(server.Clients, conn)
    server.Mutex.Unlock()

    log.Println("WebSocket connection closed")
}

func (server *WebSocketServer) HandleMessages() {
    for msg := range server.Broadcast {
        server.Mutex.Lock()
        for conn := range server.Clients {
            if err := conn.WriteJSON(msg); err != nil {
                log.Println("Error sending message:", err)
                conn.Close()
                delete(server.Clients, conn)
            }
        }
        server.Mutex.Unlock()
    }
}