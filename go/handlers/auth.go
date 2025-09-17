package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	sms "go-mongo-vue-go/libraries"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var jwtSecret = []byte("super-secret-key")
var ctx = context.Background()
var smsCodesColl *mongo.Collection

func InitSMSCollection(client *mongo.Client, dbName string) {
	smsCodesColl = client.Database(dbName).Collection("sms_codes")
	indexModel := mongo.IndexModel{
		Keys:    bson.M{"expires_at": 1},
		Options: options.Index().SetExpireAfterSeconds(0),
	}
	_, err := smsCodesColl.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		log.Println("failed to create TTL index:", err)
	}
}

type SendCodeRequest struct {
	Mobile  string `json:"mobile"`
	Country string `json:"country"`
}
type VerifyCodeRequest struct {
	Mobile  string `json:"mobile"`
	Country string `json:"country"`
	Code    string `json:"code"`
}
type LoginResponse struct {
	Token string `json:"token"`
}

func SendCode(w http.ResponseWriter, r *http.Request) {
	var req SendCodeRequest
	if smsCodesColl == nil {
		http.Error(w, "database not initialized", http.StatusInternalServerError)
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	code := fmt.Sprintf("%06d", rand.Intn(900000)+100000)
	expiresAt := time.Now().Add(5 * time.Minute)
	fullMobile := req.Country + req.Mobile
	_, err := smsCodesColl.UpdateOne(
		ctx,
		bson.M{"mobile": fullMobile},
		bson.M{"$set": bson.M{"code": code, "expires_at": expiresAt}},
		options.Update().SetUpsert(true),
	)
	if err != nil {
		http.Error(w, "failed to save code", http.StatusInternalServerError)
		return
	}
	if !(os.Getenv("SMS_SANDBOX") == "true") {
		ok, err := sms.SendSmsLogin(code, fullMobile)
		if err != nil || !ok {
			http.Error(w, "failed to send sms", http.StatusInternalServerError)
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	var result map[string]string
	if os.Getenv("SMS_SANDBOX") == "true" {
		result = map[string]string{
			"message": "code sent",
			"code":    code,
		}
	} else {
		result = map[string]string{
			"message": "code sent",
		}
	}
	json.NewEncoder(w).Encode(result)
}
func VerifyCode(w http.ResponseWriter, r *http.Request) {
	var req VerifyCodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	var result struct {
		Code      string    `bson:"code"`
		ExpiresAt time.Time `bson:"expires_at"`
	}
	fullMobile := req.Country + req.Mobile
	err := smsCodesColl.FindOne(ctx, bson.M{"mobile": fullMobile}).Decode(&result)
	if err != nil || result.Code != req.Code || time.Now().After(result.ExpiresAt) {
		http.Error(w, "invalid or expired code", http.StatusUnauthorized)
		return
	}
	_, _ = smsCodesColl.DeleteOne(ctx, bson.M{"mobile": fullMobile})
	claims := jwt.MapClaims{
		"mobile": fullMobile,
		"exp":    time.Now().Add(time.Minute * 15).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(jwtSecret)
	if err != nil {
		http.Error(w, "could not generate token", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(LoginResponse{Token: signed})
}
