package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID                primitive.ObjectID `bson:"_id,omitempty"`
	Name              string             `bson:"name" json:"name"`
	Family            string             `bson:"family" json:"family"`
	Mobile            string             `bson:"mobile" json:"mobile"`
	CodeMeli          string             `bson:"code_meli" json:"code_meli"`
	Balance           float64            `bson:"balance" json:"balance"`
	TokenBalance      float64            `bson:"token_balance" json:"token_balance"`
	FingerToken       string             `bson:"finger_token" json:"finger_token"`
	TelegramID        string             `bson:"telegram_id" json:"telegram_id"`
	Images            []string           `bson:"images" json:"images"`
	NotificationToken string             `bson:"notification_token,omitempty"`
}
