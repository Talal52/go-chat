package server

import (
    "github.com/Talal52/go-chat/chat/api"
    "log"
    "net/http"
)

func StartHTTPServer(chatHandler *api.ChatHandler, authHandler *api.AuthHandler) {
    http.HandleFunc("/messages", chatHandler.GetMessages)
    http.HandleFunc("/send", chatHandler.PostMessage)
    http.HandleFunc("/signup", authHandler.Signup)
    http.HandleFunc("/login", authHandler.Login)

    log.Println("HTTP server started on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}