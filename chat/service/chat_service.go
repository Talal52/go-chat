package service

import (
	"github.com/Talal52/go-chat/chat/db"
	"github.com/Talal52/go-chat/chat/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChatService struct {
	Repo *db.ChatRepository
}

func NewChatService(repo *db.ChatRepository) *ChatService {
	return &ChatService{Repo: repo}
}

func (s *ChatService) SaveMessage(msg models.Message) error {
	return s.Repo.SaveMessage(msg)
}

func (s *ChatService) GetMessages() ([]models.Message, error) {
	return s.Repo.GetMessages()
}
func (s *ChatService) GetMessagesByGroupID(groupID primitive.ObjectID) ([]models.Message, error) {
	return s.Repo.GetMessagesByGroupID(groupID)
}
