package server

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/Talal52/go-chat/chat/api"
	"github.com/Talal52/go-chat/chat/middleware"
)

func StartHTTPServer(chatHandler *api.ChatHandler, authHandler *api.AuthHandler) {
	router := gin.Default()

	apiGroup := router.Group("/api")
	{
		apiGroup.POST("/signup", authHandler.SignupGin)
		apiGroup.POST("/login", authHandler.LoginGin)

		// Authenticated routes
		apiGroup.Use(middleware.AuthMiddleware())
		{
			apiGroup.GET("/messages", chatHandler.GetMessagesGin)
			apiGroup.POST("/send-message", chatHandler.PostMessageGin) // ✅ Fixed here
		}
	}

	log.Println("HTTP server started on :8080")
	router.Run(":8080")
}
