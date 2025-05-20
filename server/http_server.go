package server

import (
	"log"
	"net/http"

	"github.com/Talal52/go-chat/chat/api"
	"github.com/Talal52/go-chat/chat/middleware"
	"github.com/gin-gonic/gin"
)

func StartHTTPServer(chatHandler *api.ChatHandler, authHandler *api.AuthHandler) {
	router := gin.Default()

	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	apiGroup := router.Group("/api")
	{
		apiGroup.POST("/signup", authHandler.SignupGin)
		apiGroup.POST("/login", authHandler.LoginGin)

		// Authenticated routes
		apiGroup.Use(middleware.AuthMiddleware())
		{
			apiGroup.GET("/messages", chatHandler.GetMessagesGin)
			apiGroup.POST("/send-message", chatHandler.PostMessageGin)
		}
	}

	log.Println("HTTP server started on :8080")
	router.Run(":8080")
}
