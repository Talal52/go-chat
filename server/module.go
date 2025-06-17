package server

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/Talal52/go-chat/chat/api"
	"github.com/Talal52/go-chat/chat/db"
	"github.com/Talal52/go-chat/chat/service"
	"github.com/Talal52/go-chat/server/websocket" // Import the websocket package
	"go.mongodb.org/mongo-driver/mongo"
)

func InitServers(mongoDB *mongo.Database, postgresDB *sql.DB) {
	userRepo := db.NewUserRepository(postgresDB)
	chatRepo := db.NewChatRepository(mongoDB)

	authService := service.NewAuthService(userRepo)
	chatService := service.NewChatService(chatRepo)

	authHandler := api.NewAuthHandler(authService)
	chatHandler := api.NewChatHandler(chatService)

	// Start HTTP server
	go StartHTTPServer(chatHandler, authHandler)

	// Start WebSocket server
	webSocketServer := websocket.NewWebSocketServer(chatService) // Use websocket package
	go webSocketServer.HandleMessages()

	// Configure WebSocket endpoint
	http.HandleFunc("/ws", webSocketServer.HandleConnections)

	log.Println("WebSocket server started on :8081/ws")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
