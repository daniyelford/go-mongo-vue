package handlers

import (
	"encoding/json"
	"fmt"
	"go-mongo-vue-users/config"
	"go-mongo-vue-users/middleware"
	"go-mongo-vue-users/models"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func RefreshToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var body struct {
		RefreshToken string `json:"refreshToken"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, fmt.Sprintf("invalid request body: %v", err), http.StatusBadRequest)
		return
	}
	if body.RefreshToken == "" {
		http.Error(w, "missing refresh token", http.StatusBadRequest)
		return
	}
	token, err := jwt.Parse(body.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return config.JwtSecret, nil
	})
	if err != nil || !token.Valid {
		http.Error(w, "invalid refresh token", http.StatusUnauthorized)
		return
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		http.Error(w, "invalid token claims", http.StatusUnauthorized)
		return
	}
	mobile, ok := claims["mobile"].(string)
	if !ok || mobile == "" {
		http.Error(w, "invalid refresh token", http.StatusUnauthorized)
		return
	}
	val, err := config.RedisClient.Get(config.Ctx, "refresh:"+mobile).Result()
	if err != nil || val != body.RefreshToken {
		http.Error(w, "refresh token expired", http.StatusUnauthorized)
		return
	}
	newClaims := jwt.MapClaims{
		"mobile": mobile,
		"exp":    time.Now().Add(15 * time.Minute).Unix(),
	}
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)
	signed, err := newToken.SignedString(config.JwtSecret)
	if err != nil {
		http.Error(w, "cannot generate new access token", http.StatusInternalServerError)
		return
	}
	config.RedisClient.Set(config.Ctx, "token:"+mobile, signed, 15*time.Minute)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"accessToken": signed})
}
func ValidateToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	mobile := r.Context().Value(middleware.MobileContextKey).(string)
	userColl := config.MongoClient.Database(os.Getenv("DB_NAME")).Collection("users")
	var user models.User
	userHasInfo := true
	if err := userColl.FindOne(config.Ctx, bson.M{"mobile": mobile}).Decode(&user); err == mongo.ErrNoDocuments {
		userHasInfo = false
	} else if err != nil {
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]any{
		"loggedIn":    true,
		"userHasInfo": userHasInfo,
		"expired":     false,
	})
}
