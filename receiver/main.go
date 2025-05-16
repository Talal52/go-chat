package main

import (
    "fmt"
    "log"
    "github.com/Talal52/go-chat/config"
    "github.com/Talal52/go-chat/server/websocket"
)

func main() {
    cfg := config.LoadConfig()
    wsClient := websocket.NewClient(cfg.WebSocketURL)
    
    err := wsClient.Connect()
    if err != nil {
        log.Fatalf("Failed to connect to websocket server: %v", err)
    }

    fmt.Println("Connected to WebSocket server as receiver")
    wsClient.ListenMessages()
}
