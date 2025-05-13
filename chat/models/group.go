package models

import (
    "go.mongodb.org/mongo-driver/bson/primitive"
    "time"
)

type Group struct {
    ID        primitive.ObjectID `bson:"_id,omitempty"`
    Name      string             `bson:"name"`
    Members   []string           `bson:"members"`
    CreatedAt time.Time          `bson:"createdAt"`
}