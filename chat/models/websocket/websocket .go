package websocket

import (
	"sync"

	"github.com/Talal52/go-chat/chat/service"
	"github.com/gorilla/websocket"
)

type WebSocketServer struct {
	Clients   map[*websocket.Conn]string // Map of WebSocket connections to usernames
	Broadcast chan interface{}           // Channel for broadcasting messages
	Mutex     sync.Mutex                 // Mutex to protect the Clients map
	Service   *service.ChatService       // Chat service for saving messages
}
