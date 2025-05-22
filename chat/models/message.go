package models

import (
	"time"
)

type Message struct {
	ID         string    `json:"id" bson:"_id,omitempty"`
	SenderID   string    `json:"senderId" bson:"sender_id"`
	ReceiverID string    `json:"receiverId" bson:"receiver_id"`
	Message    string    `json:"message" bson:"message"`
	Timestamp  time.Time `json:"timestamp" bson:"timestamp"`
}
