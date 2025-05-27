package service

import (
	"context"

	"github.com/Talal52/go-chat/chat/db"
	"github.com/Talal52/go-chat/chat/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChatService struct {
	Repo *db.ChatRepository
}

func NewChatService(repo *db.ChatRepository) *ChatService {
	return &ChatService{Repo: repo}
}

func (s *ChatService) GetMessages() ([]models.Message, error) {
	// Fetch messages
	ctx := context.TODO()
	cursor, err := s.Repo.Collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var messages []models.Message
	if err := cursor.All(ctx, &messages); err != nil {
		return nil, err
	}
	return messages, nil
}

func (s *ChatService) SaveMessage(msg models.Message) error {
	// Save the message to MongoDB
	ctx := context.TODO()
	_, err := s.Repo.Collection.InsertOne(ctx, msg)
	return err
}

func (s *ChatService) GetMessagesByGroupID(groupID primitive.ObjectID) ([]models.Message, error) {
	return s.Repo.GetMessagesByGroupID(groupID)
}
