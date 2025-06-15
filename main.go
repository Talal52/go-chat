package main

import (
	"log"
	"net/http"

	"github.com/Talal52/go-chat/chat/api"
	"github.com/Talal52/go-chat/chat/db"
	"github.com/Talal52/go-chat/chat/service"
	"github.com/Talal52/go-chat/config"
	"github.com/Talal52/go-chat/server/websocket"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors" // Add CORS import
	// "go.mongodb.org/mongo-driver/mongo"
)

func main() {
	// Connect to databases
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
	authHandler := api.NewAuthHandler(authService)
	chatHandler := api.NewChatHandler(chatService)

	// Set up Gin router
	router := gin.Default()

	// Configure CORS
	config := cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}
	router.Use(cors.New(config))

	// Register auth endpoints (public)
	router.POST("/api/signup", authHandler.SignupGin)
	router.POST("/api/login", authHandler.LoginGin)

	// Register chat endpoints (protected)
	chatRoutes := router.Group("/api")
	chatRoutes.Use(api.AuthMiddleware())
	{
		chatRoutes.GET("/messages", chatHandler.GetMessagesGin)
		chatRoutes.POST("/send-message", chatHandler.PostMessageGin)
		chatRoutes.GET("/users", chatHandler.GetUsers)
	}

	// Start HTTP server
	go func() {
		log.Println("Starting HTTP server on :8080...")
		if err := router.Run(":8080"); err != nil {
			log.Fatal("HTTP server failed:", err)
		}
	}()

	// Start WebSocket server
	wsServer := websocket.NewWebSocketServer(chatService)
	go wsServer.HandleMessages()

	log.Println("Starting WebSocket server on :8083...")
	http.HandleFunc("/ws", wsServer.HandleConnections)
	if err := http.ListenAndServe(":8083", nil); err != nil {
		log.Fatal("WebSocket server failed:", err)
	}
}