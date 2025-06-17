package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/Talal52/go-chat/chat/api"
	"github.com/Talal52/go-chat/chat/db"
	"github.com/Talal52/go-chat/chat/service"
	"github.com/Talal52/go-chat/config"
	"github.com/Talal52/go-chat/server"
)

func main() {
	cfg := config.LoadConfig()

	// Connect to databases
	postgresDB := config.ConnectPostgres(cfg.PostgresURI)
	defer postgresDB.Close()

	mongoDB := config.ConnectDB(cfg.MongoURI, cfg.DBName)
	defer mongoDB.Client().Disconnect(context.Background())

	// Initialize repositories
	userRepo := db.NewUserRepository(postgresDB)
	// chatRepo := db.NewChatRepository(mongoDB) // Uncomment if needed

	// Initialize services
	authService := service.NewAuthService(userRepo)
	// chatService := service.NewChatService(chatRepo) // Uncomment if needed

	// Create authHandler
	authHandler := api.NewAuthHandler(authService)

	// Start HTTP server
	srv := server.NewHTTPServer(cfg, authHandler)
	fmt.Printf("HTTP Sender running on port %s\n", cfg.HTTPPort)
	log.Fatal(http.ListenAndServe(":"+cfg.HTTPPort, srv))
}
