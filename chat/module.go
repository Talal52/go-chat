package chat

import (
	"database/sql"
	"github.com/Talal52/go-chat/chat/api"
	"github.com/Talal52/go-chat/chat/db"
	"github.com/Talal52/go-chat/chat/service"
)

func InitChatModule(dbConn *sql.DB) *api.ChatHandler {
	repo := db.NewChatRepository(dbConn)
	svc := service.NewChatService(repo)
	return api.NewChatHandler(svc)
}
