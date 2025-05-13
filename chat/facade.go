package chat

import (
	"github.com/Talal52/go-chat/chat/api"
	chat_db "github.com/Talal52/go-chat/chat/db"
	"github.com/Talal52/go-chat/chat/service"
	"go.mongodb.org/mongo-driver/mongo"
)

type ChatFacade struct {
	Handler *api.ChatHandler
	Service *service.ChatService
}

func NewChatFacade(db *mongo.Database) *ChatFacade {
	repo := chat_db.NewChatRepository(db)

	svc := service.NewChatService(repo)

	handler := api.NewChatHandler(svc)

	return &ChatFacade{
		Handler: handler,
		Service: svc,
	}
}
