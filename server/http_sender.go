package server

import (
	"github.com/gin-gonic/gin"
	"github.com/Talal52/go-chat/config"
)

func NewHTTPServer(cfg *config.Config) *gin.Engine {
	router := gin.Default()

	router.POST("/send", func(c *gin.Context) {
		// You can send a message to websocket server from here using HTTP
		c.JSON(200, gin.H{"status": "message sent!"})
	})

	return router
}
