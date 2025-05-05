package models

import "time"

type Message struct {
    ID        string    `bson:"_id,omitempty" json:"id"`
    Username  string    `bson:"username" json:"username"`
    Content   string    `bson:"content" json:"content"`
    CreatedAt time.Time `bson:"created_at" json:"created_at"`
}