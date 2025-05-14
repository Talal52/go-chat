package main

import (
    "log"
    "net/http"

    "github.com/Talal52/go-chat/chat/api"
    "github.com/Talal52/go-chat/chat/service"
    "github.com/Talal52/go-chat/server"
    "github.com/Talal52/go-chat/server/websocket"
)

func main() {
    // Initialize services
    chatService := service.NewChatService()
    authService := service.NewAuthService() // Ensure AuthService is initialized

    // Initialize handlers
    chatHandler := &api.ChatHandler{Service: chatService}
    authHandler := &api.AuthHandler{Service: authService}

    // Start HTTP server
    go func() {
        log.Println("Starting HTTP server...")
        server.StartHTTPServer(chatHandler, authHandler)
    }()

    // Start WebSocket server
    wsServer := websocket.NewWebSocketModule(chatService)
    go wsServer.HandleMessages()

    log.Println("Starting WebSocket server on :8081...")
    http.HandleFunc("/ws", websocket.WebSocketHandler(wsServer))
    log.Fatal(http.ListenAndServe(":8081", nil))
}