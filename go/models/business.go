package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Media struct {
	URL      string `bson:"url" json:"url"`
	Type     string `bson:"type" json:"type"`
	Filename string `bson:"filename" json:"filename"`
}
type Business struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID      primitive.ObjectID `bson:"user_id" json:"user_id"`
	Name        string             `bson:"name" json:"name"`
	Slug        string             `bson:"slug" json:"slug"`
	Description string             `bson:"description" json:"description"`
	Category    string             `bson:"category" json:"category"`
	Phone       string             `bson:"phone" json:"phone"`
	Address     string             `bson:"address" json:"address"`
	Location    GeoPoint           `bson:"location" json:"location"`
	Media       []Media            `bson:"media" json:"media"`
	Status      string             `bson:"status" json:"status"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}

type GeoPoint struct {
	Type        string    `bson:"type" json:"type"`
	Coordinates []float64 `bson:"coordinates" json:"coordinates"` // [longitude, latitude]
}
