package websocket

import (
	"sync"

	"github.com/Talal52/go-chat/chat/models"
	"github.com/Talal52/go-chat/chat/service"
	"github.com/gorilla/websocket"
)

type WebSocketServer struct {
	Clients   map[*websocket.Conn]string 
	Broadcast chan models.Message          
	Mutex     sync.Mutex                 
	Service   *service.ChatService      
}
