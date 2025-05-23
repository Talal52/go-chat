package chat

import (
	"github.com/Talal52/go-chat/chat/api"
	chat_db "github.com/Talal52/go-chat/chat/db"
	"github.com/Talal52/go-chat/chat/service"
	"go.mongodb.org/mongo-driver/mongo"
)

func InitChatModule(db *mongo.Database) *api.ChatHandler {
	repo := chat_db.NewChatRepository(db)

	svc := service.NewChatService(repo)

	return api.NewChatHandler(svc)
}
