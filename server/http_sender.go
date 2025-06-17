package server

import (
	"github.com/Talal52/go-chat/chat/api"
	"github.com/Talal52/go-chat/config"
	"github.com/gin-gonic/gin"
)

func NewHTTPServer(cfg *config.Config, authHandler *api.AuthHandler) *gin.Engine {
	router := gin.Default()

	// Auth-related routes
	if authHandler != nil {
		router.POST("/signup", authHandler.SignupGin)
		router.POST("/login", authHandler.LoginGin)
	}

	// Optional: Keep or remove the /send route based on your needs
	router.POST("/send", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "message sent!"})
	})

	return router
}
