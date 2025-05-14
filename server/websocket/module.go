package websocket

import "github.com/Talal52/go-chat/chat/service"

func NewWebSocketModule(chatService *service.ChatService) *WebSocketServer {
    return NewWebSocketServer(chatService)
}