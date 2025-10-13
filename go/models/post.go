package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PostMedia struct {
	URL      string `bson:"url" json:"url"`
	Type     string `bson:"type" json:"type"`
	Filename string `bson:"filename" json:"filename"`
}
type Post struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	UserID    primitive.ObjectID `bson:"user_id"`
	Title     string             `bson:"title"`
	Media     []PostMedia        `bson:"media" json:"media"`
	Content   string             `bson:"content"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}
