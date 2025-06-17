package main

import (
	"context"
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
	cfg := config.LoadConfig()

	postgresDB := config.ConnectPostgres(cfg.PostgresURI)
	defer postgresDB.Close()

	mongoDB := config.ConnectDB(cfg.MongoURI, cfg.DBName)
	defer mongoDB.Client().Disconnect(context.Background())

	userRepo := db.NewUserRepository(postgresDB)
	chatRepo := db.NewChatRepository(mongoDB)

	authService := service.NewAuthService(userRepo)
	chatService := service.NewChatService(chatRepo)

	authHandler := api.NewAuthHandler(authService)

	go func() {
		log.Println("Starting HTTP server on :8080...")
		router := server.NewHTTPServer(cfg, authHandler)
		if err := router.Run(":8080"); err != nil {
			log.Fatal("HTTP server failed:", err)
		}
	}()

	wsServer := websocket.NewWebSocketModule(chatService)
	go wsServer.HandleMessages()

	log.Println("Starting WebSocket server on :8081...")
	http.HandleFunc("/ws", websocket.WebSocketHandler(wsServer))
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatal("WebSocket server failed:", err)
	}
}