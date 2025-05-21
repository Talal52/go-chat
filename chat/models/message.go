package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Message struct {
	ID         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Sender     string             `json:"sender" bson:"sender"`
	Content    string             `json:"message" bson:"content"`
	ReceiverId int                `json:"receiver_id" bson:"receiver_id"`
	GroupID    primitive.ObjectID `json:"group_id,omitempty" bson:"group_id,omitempty"`
	CreatedAt  time.Time          `json:"created_at" bson:"created_at"`
}
