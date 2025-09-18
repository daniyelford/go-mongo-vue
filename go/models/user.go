package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type WebAuthnCredential struct {
	CredentialID []byte `bson:"credential_id" json:"credential_id"`
	PublicKey    []byte `bson:"public_key" json:"public_key"`
	SignCount    uint32 `bson:"sign_count" json:"sign_count"`
}
type User struct {
	ID                primitive.ObjectID   `bson:"_id,omitempty"`
	Name              string               `bson:"name" json:"name"`
	Family            string               `bson:"family" json:"family"`
	Mobile            string               `bson:"mobile" json:"mobile"`
	CodeMeli          string               `bson:"code_meli" json:"code_meli"`
	Balance           float64              `bson:"balance" json:"balance"`
	TokenBalance      float64              `bson:"token_balance" json:"token_balance"`
	FingerTokens      []WebAuthnCredential `bson:"finger_tokens" json:"finger_tokens"`
	TelegramID        string               `bson:"telegram_id" json:"telegram_id"`
	Image             string               `bson:"image" json:"image"`
	NotificationToken string               `bson:"notification_token,omitempty"`
}
