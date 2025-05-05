package chat

import (
    "go.mongodb.org/mongo-driver/mongo"
    "github.com/Talal52/go-chat/chat/api"
    chat_db "github.com/Talal52/go-chat/chat/db"
    "github.com/Talal52/go-chat/chat/service"
)

type ChatFacade struct {
    Handler *api.ChatHandler
    Service *service.ChatService
}

func NewChatFacade(db *mongo.Database) *ChatFacade {
    // Initialize repository
    repo := chat_db.NewChatRepository(db)

    // Initialize service
    svc := service.NewChatService(repo)

    // Initialize API handler
    handler := api.NewChatHandler(svc)

    return &ChatFacade{
        Handler: handler,
        Service: svc,
    }
}