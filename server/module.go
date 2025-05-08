package server

import (
	"log"
	"net/http"

	"database/sql"

	"github.com/Talal52/go-chat/chat"
	"go.mongodb.org/mongo-driver/mongo"
)

func InitServers(mongoDB *mongo.Database, postgresDB *sql.DB) {
	// Initialize Chat Module and WebSocket Server
	chatHandler := chat.InitChatModule(mongoDB)
	webSocketServer := NewWebSocketServer(chatHandler.Service)

	// Start WebSocket message handling in a separate goroutine
	go webSocketServer.HandleMessages()

	// Configure and start HTTP server for WebSocket connections
	http.HandleFunc("/ws", webSocketServer.HandleConnections)
	log.Println("WebSocket server started on :8080/ws")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
}
