package db

import (
    "context"
    "time"

    "github.com/Talal52/go-chat/chat/models"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
)

type ChatRepository struct {
    Collection *mongo.Collection
}

func NewChatRepository(db *mongo.Database) *ChatRepository {
    return &ChatRepository{
        Collection: db.Collection("messages"),
    }
}

func (r *ChatRepository) SaveMessage(msg models.Message) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    _, err := r.Collection.InsertOne(ctx, msg)
    return err
}

func (r *ChatRepository) GetMessages() ([]models.Message, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    cursor, err := r.Collection.Find(ctx, bson.M{})
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)

    var messages []models.Message
    for cursor.Next(ctx) {
        var msg models.Message
        if err := cursor.Decode(&msg); err != nil {
            return nil, err
        }
        messages = append(messages, msg)
    }

    return messages, nil
}