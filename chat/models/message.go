package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Message struct {
	Sender    string              `json:"sender"`
	Receiver  string              `json:"receiver"`
	Content   string              `json:"content"`
	CreatedAt time.Time           `json:"created_at"`
	GroupID   *primitive.ObjectID `json:"group_id,omitempty" bson:"group_id,omitempty"`
}
