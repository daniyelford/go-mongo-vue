package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Reservation struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ServiceID  primitive.ObjectID `bson:"service_id" json:"service_id"`
	BusinessID primitive.ObjectID `bson:"business_id" json:"business_id"`
	UserID     primitive.ObjectID `bson:"user_id" json:"user_id"`
	TimeSlot   TimeSlot           `bson:"time_slot" json:"time_slot"`
	Status     string             `bson:"status" json:"status"`
	Note       string             `bson:"note,omitempty" json:"note,omitempty"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at" json:"updated_at"`
}

type TimeSlot struct {
	Start time.Time `bson:"start" json:"start"`
	End   time.Time `bson:"end" json:"end"`
}
