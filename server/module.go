package server

import (
    "github.com/Talal52/go-chat/chat"
    "github.com/Talal52/go-chat/chat/api"
    "github.com/Talal52/go-chat/chat/db"
    "github.com/Talal52/go-chat/chat/service"
    "go.mongodb.org/mongo-driver/mongo"
    "database/sql"
)

func InitServers(mongoDB *mongo.Database, postgresDB *sql.DB) {
    // Initialize Chat Module
    chatHandler := chat.InitChatModule(mongoDB)

    // Initialize User Module
    userRepo := db.NewUserRepository(postgresDB)
    authService := service.NewAuthService(userRepo)
    authHandler := api.NewAuthHandler(authService)

    // Start Servers
    go StartHTTPServer(chatHandler, authHandler) // Pass both handlers
    go StartTCPServer(chatHandler.Service)
}