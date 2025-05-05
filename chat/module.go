package chat

import (
	"github.com/Talal52/go-chat/chat/api"
	chat_db "github.com/Talal52/go-chat/chat/db"
	"github.com/Talal52/go-chat/chat/service"
	"go.mongodb.org/mongo-driver/mongo"
)

func InitChatModule(db *mongo.Database) *api.ChatHandler {
	// Initialize the repository
	repo := chat_db.NewChatRepository(db)

	// Initialize the service
	svc := service.NewChatService(repo)

	// Return the API handler
	return api.NewChatHandler(svc)
}
