package server

import (
	"log"
	"net/http"

	"github.com/Talal52/go-chat/chat/api"
	"github.com/Talal52/go-chat/chat/db"
	"github.com/Talal52/go-chat/chat/service"
	"github.com/Talal52/go-chat/server/websocket"
	"github.com/gin-gonic/gin"
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

	router := gin.Default()

	router.POST("/api/signup", authHandler.SignupGin)
	router.POST("/api/login", authHandler.LoginGin)

	chatRoutes := router.Group("/api")
	chatRoutes.Use(api.AuthMiddleware())
	{
		chatRoutes.GET("/messages", chatHandler.GetMessagesGin)
		chatRoutes.POST("/send-message", chatHandler.PostMessageGin)
	}

	go func() {
		log.Println("Starting HTTP server on :8080...")
		if err := router.Run(":8080"); err != nil {
			log.Fatal("HTTP server failed:", err)
		}
	}()

	wsServer := websocket.NewWebSocketServer(chatService)
	go wsServer.HandleMessages()

	log.Println("Starting WebSocket server on :8081...")
	http.HandleFunc("/ws", wsServer.HandleConnections)
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatal("WebSocket server failed:", err)
	}
}