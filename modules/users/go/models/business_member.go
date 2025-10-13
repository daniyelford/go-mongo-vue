package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BusinessMember struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	BusinessID  primitive.ObjectID `bson:"business_id" json:"business_id"`
	UserID      primitive.ObjectID `bson:"user_id" json:"user_id"`
	Role        string             `bson:"role" json:"role"` // owner | manager | exp
	Permissions Permissions        `bson:"permissions" json:"permissions"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}

type Permissions struct {
	CanView               bool `bson:"can_view" json:"can_view"`
	CanEditServices       bool `bson:"can_edit_services" json:"can_edit_services"`
	CanManageReservations bool `bson:"can_manage_reservations" json:"can_manage_reservations"`
}
