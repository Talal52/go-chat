package server

import (
	"database/sql"
	"github.com/Talal52/go-chat/chat"
)

func InitServers(dbConn *sql.DB) {
	handler := chat.InitChatModule(dbConn)
	go StartHTTPServer(handler)
	go StartTCPServer() // You can pass ChatHandler or ChatService
}
