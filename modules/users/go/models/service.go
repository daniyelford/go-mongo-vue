package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	BusinessID  primitive.ObjectID `bson:"business_id" json:"business_id"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description" json:"description"`
	Duration    int                `bson:"duration" json:"duration"` // in minutes
	Price       int64              `bson:"price" json:"price"`
	Available   bool               `bson:"available" json:"available"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}
