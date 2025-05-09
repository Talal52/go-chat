package server

import (
    "github.com/Talal52/go-chat/chat/api"
    "log"
    "net/http"
)

func StartHTTPServer(chatHandler *api.ChatHandler, authHandler *api.AuthHandler) {
    routes := map[string]http.HandlerFunc{
        "/messages": chatHandler.GetMessages,
        "/send":     chatHandler.PostMessage,
        "/signup":   authHandler.Signup, 
        "/login":    authHandler.Login, 
    }

    for route, handler := range routes {
        http.HandleFunc(route, handler)
    }

    log.Println("HTTP server started on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}