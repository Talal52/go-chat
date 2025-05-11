package server

import (
    "log"

    "github.com/gin-gonic/gin"
    "github.com/Talal52/go-chat/chat/api"
)

func StartHTTPServer(chatHandler *api.ChatHandler, authHandler *api.AuthHandler) {
    router := gin.Default()

    // ✅ Serve static files under /static instead of /
    router.Static("/static", "./frontend")

    // ✅ Group API routes under /api
    apiGroup := router.Group("/api")
    {
        apiGroup.POST("/signup", authHandler.SignupGin)
        apiGroup.POST("/login", authHandler.LoginGin)
        apiGroup.GET("/messages", chatHandler.GetMessagesGin)
    }

    log.Println("HTTP server started on :8080")
    router.Run(":8080")
}
