package models

import "time"

type Message struct {
    ID        string    `bson:"_id,omitempty" json:"id"`       // MongoDB Object ID
    Sender    string    `bson:"sender" json:"sender"`         // Sender's username
    Receiver  string    `bson:"receiver" json:"receiver"`     // Receiver's username
    Content   string    `bson:"content" json:"content"`       // Message content
    CreatedAt time.Time `bson:"created_at" json:"created_at"` // Timestamp
}