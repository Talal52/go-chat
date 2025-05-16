package main

import (
    "log"
    "net/http"

    "github.com/Talal52/go-chat/chat/api"
    "github.com/Talal52/go-chat/chat/db"
    "github.com/Talal52/go-chat/chat/service"
    "github.com/Talal52/go-chat/config"
    "github.com/Talal52/go-chat/server"
    "github.com/Talal52/go-chat/server/websocket"
)

func main() {
    postgresDB := config.ConnectPostgres()
    defer postgresDB.Close()

    mongoDB := config.ConnectDB()

    // Initialize repositories
    userRepo := db.NewUserRepository(postgresDB)
    chatRepo := db.NewChatRepository(mongoDB)

    // Initialize services
    authService := service.NewAuthService(userRepo)
    chatService := service.NewChatService(chatRepo)

    // Initialize handlers
    chatHandler := &api.ChatHandler{Service: chatService}
    authHandler := &api.AuthHandler{Service: authService}

    go func() {
        log.Println("Starting HTTP server...")
        server.StartHTTPServer(chatHandler, authHandler)
    }()

    wsServer := websocket.NewWebSocketModule(chatService)
    go wsServer.HandleMessages()

    log.Println("Starting WebSocket server on :8081...")
    http.HandleFunc("/ws", websocket.WebSocketHandler(wsServer))
    log.Fatal(http.ListenAndServe(":8081", nil))
}