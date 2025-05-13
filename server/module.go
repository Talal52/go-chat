package server

import (
    "log"
    "net/http"

    "github.com/Talal52/go-chat/chat/api"
    "github.com/Talal52/go-chat/chat/db"
    "github.com/Talal52/go-chat/chat/service"
    "go.mongodb.org/mongo-driver/mongo"
    "database/sql"
)

func InitServers(mongoDB *mongo.Database, postgresDB *sql.DB) {
    userRepo := db.NewUserRepository(postgresDB)
    chatRepo := db.NewChatRepository(mongoDB)

    authService := service.NewAuthService(userRepo)
    chatService := service.NewChatService(chatRepo)

    authHandler := api.NewAuthHandler(authService)
    chatHandler := api.NewChatHandler(chatService)

    go StartHTTPServer(chatHandler, authHandler)

    // Start WebSocket server
    webSocketServer := NewWebSocketServer(chatService)
    go webSocketServer.HandleMessages()

    // Configure WebSocket endpoint
    http.HandleFunc("/ws", webSocketServer.HandleConnections)

	log.Println("WebSocket server started on :8081/ws")
	log.Fatal(http.ListenAndServe(":8081", nil))
}