package server

import (
    "log"

    "github.com/gin-gonic/gin"
    "github.com/Talal52/go-chat/chat/api"
)

func StartHTTPServer(chatHandler *api.ChatHandler, authHandler *api.AuthHandler) {
    router := gin.Default()

    router.Static("/static", "./frontend")

    apiGroup := router.Group("/api")
    {
        apiGroup.POST("/signup", authHandler.SignupGin)
        apiGroup.POST("/login", authHandler.LoginGin)
        apiGroup.GET("/messages", chatHandler.GetMessagesGin)
        apiGroup.GET("/group/messages?group_id=",chatHandler.GetGroupMessages)
    }

    log.Println("HTTP server started on :8080")
    router.Run(":8080")
}
