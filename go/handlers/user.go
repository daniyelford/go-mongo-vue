package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	sms "go-mongo-vue-go/libraries"
	"go-mongo-vue-go/models"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))
var ctx = context.Background()
var redisClient *redis.Client
var mongoClient *mongo.Client

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
	Token   string `json:"token"`
	NewUser bool   `json:"newUser"`
}

func InitRedis(client *redis.Client) {
	redisClient = client
}
func InitMongo(client *mongo.Client) {
	mongoClient = client
}
func SendCode(w http.ResponseWriter, r *http.Request) {
	var req SendCodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	code := fmt.Sprintf("%06d", rand.Intn(900000)+100000)
	fullMobile := req.Country + req.Mobile
	err := redisClient.Set(ctx, "sms:"+fullMobile, code, 5*time.Minute).Err()
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
	fullMobile := req.Country + req.Mobile
	storedCode, err := redisClient.Get(ctx, "sms:"+fullMobile).Result()
	if err == redis.Nil || storedCode != req.Code {
		http.Error(w, "invalid or expired code", http.StatusUnauthorized)
		return
	}
	redisClient.Del(ctx, "sms:"+fullMobile)
	var user models.User
	userColl := mongoClient.Database(os.Getenv("DB_NAME")).Collection("users")
	err = userColl.FindOne(ctx, bson.M{"mobile": fullMobile}).Decode(&user)
	newUser := false
	if err == mongo.ErrNoDocuments {
		newUser = true
	} else if err != nil {
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}
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
	err = redisClient.Set(ctx, "token:"+fullMobile, signed, 25000*time.Minute).Err()
	if err != nil {
		http.Error(w, "cannot store token", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(LoginResponse{Token: signed, NewUser: newUser})
}
func ValidateToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		json.NewEncoder(w).Encode(map[string]any{
			"loggedIn": false,
		})
		return
	}
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		json.NewEncoder(w).Encode(map[string]any{
			"loggedIn": false,
		})
		return
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		json.NewEncoder(w).Encode(map[string]any{
			"loggedIn": false,
		})
		return
	}
	mobile, ok := claims["mobile"].(string)
	if !ok || mobile == "" {
		json.NewEncoder(w).Encode(map[string]any{
			"loggedIn": false,
		})
		return
	}
	val, err := redisClient.Get(ctx, "token:"+mobile).Result()
	if err != nil || val != tokenString {
		json.NewEncoder(w).Encode(map[string]any{
			"loggedIn": false,
		})
		return
	}
	var user models.User
	userColl := mongoClient.Database(os.Getenv("DB_NAME")).Collection("users")
	err = userColl.FindOne(ctx, bson.M{"mobile": mobile}).Decode(&user)
	userHasInfo := true
	if err == mongo.ErrNoDocuments {
		userHasInfo = false
	} else if err != nil {
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]any{
		"loggedIn":    true,
		"userHasInfo": userHasInfo,
	})
}
func Logout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		json.NewEncoder(w).Encode(map[string]any{
			"loggedIn": false,
		})
		return
	}
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		json.NewEncoder(w).Encode(map[string]any{
			"loggedIn": false,
		})
		return
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		json.NewEncoder(w).Encode(map[string]any{
			"loggedIn": false,
		})
		return
	}
	mobile, ok := claims["mobile"].(string)
	if !ok || mobile == "" {
		json.NewEncoder(w).Encode(map[string]any{
			"loggedIn": false,
		})
		return
	}
	val, err := redisClient.Get(ctx, "token:"+mobile).Result()
	if err != nil || val != tokenString {
		json.NewEncoder(w).Encode(map[string]any{
			"loggedIn": false,
		})
		return
	}
	if _, err := redisClient.Get(ctx, "token:"+mobile).Result(); err == redis.Nil {
		http.Error(w, "invalid or expired code", http.StatusUnauthorized)
		return
	}
	redisClient.Del(ctx, "token:"+mobile)
	json.NewEncoder(w).Encode(map[string]any{
		"loggedIn": false,
		"Logout":   true,
	})
}
