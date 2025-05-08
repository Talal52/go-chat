package server

import (
	"log"
	"net/http"

	"database/sql"

	"github.com/Talal52/go-chat/chat"
	"github.com/Talal52/go-chat/chat/models/websocket"
	"go.mongodb.org/mongo-driver/mongo"
)

func InitServers(mongoDB *mongo.Database, postgresDB *sql.DB) {
	// Initialize Chat Module
	chatHandler := chat.InitChatModule(mongoDB)

	// Initialize WebSocket Server
	webSocketServer := NewWebSocketServer(&websocket.WebSocketServer{
		Service: chatHandler.Service,
	})
	go webSocketServer.HandleMessages()

	// Start HTTP Server
	http.HandleFunc("/ws", webSocketServer.HandleConnections) // WebSocket endpoint
	log.Println("WebSocket server started on :8080/ws")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
