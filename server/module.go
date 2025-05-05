package server

import (
    "github.com/Talal52/go-chat/chat"
    "go.mongodb.org/mongo-driver/mongo"
)

func InitServers(dbConn *mongo.Database) {
    // Initialize the Chat Module
    handler := chat.InitChatModule(dbConn)

    // Start Servers
    go StartHTTPServer(handler)
    go StartTCPServer(handler.Service)
}