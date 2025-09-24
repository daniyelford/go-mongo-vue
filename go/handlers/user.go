package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	sms "go-mongo-vue-go/libraries"
	"go-mongo-vue-go/models"
	"go-mongo-vue-go/service"
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

var (
	jwtSecret   = []byte(os.Getenv("JWT_SECRET"))
	ctx         = context.Background()
	redisClient *redis.Client
	mongoClient *mongo.Client
)

// Request structs
type SendCodeRequest struct {
	Mobile  string `json:"mobile"`
	Country string `json:"country"`
}

type VerifyCodeRequest struct {
	Mobile  string `json:"mobile"`
	Country string `json:"country"`
	Code    string `json:"code"`
}

// Initialize Redis and Mongo clients
func InitRedis(client *redis.Client) { redisClient = client }
func InitMongo(client *mongo.Client) { mongoClient = client }

// Send SMS code
func SendCode(w http.ResponseWriter, r *http.Request) {
	var req SendCodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	fullMobile := req.Country + req.Mobile
	limitKey := "sms_limit:" + fullMobile
	lastKey := "sms_last:" + fullMobile

	// Check last send timestamp (15 sec delay)
	lastTs, _ := redisClient.Get(ctx, lastKey).Int64()
	if time.Now().Unix()-lastTs < 15 {
		http.Error(w, "please wait 15 seconds before retrying", http.StatusTooManyRequests)
		return
	}

	// Check max 3 sends in 5 minutes
	count, _ := redisClient.Get(ctx, limitKey).Int()
	if count >= 3 {
		http.Error(w, "too many requests, try again in 5 minutes", http.StatusTooManyRequests)
		return
	}

	// Generate code
	code := fmt.Sprintf("%06d", rand.Intn(900000)+100000)
	if err := redisClient.Set(ctx, "sms:"+fullMobile, code, 5*time.Minute).Err(); err != nil {
		http.Error(w, "failed to save code", http.StatusInternalServerError)
		return
	}

	if os.Getenv("SMS_SANDBOX") != "true" {
		ok, err := sms.SendSmsLogin(code, fullMobile)
		if err != nil || !ok {
			http.Error(w, "failed to send sms", http.StatusInternalServerError)
			return
		}
	}

	// Increment counter and update last send time
	redisClient.Incr(ctx, limitKey)
	redisClient.Expire(ctx, limitKey, 5*time.Minute)
	redisClient.Set(ctx, lastKey, time.Now().Unix(), 5*time.Minute)

	// Return result
	result := map[string]string{"message": "code sent"}
	if os.Getenv("SMS_SANDBOX") == "true" {
		result["code"] = code
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// Verify code with 3 attempts max
func VerifyCode(w http.ResponseWriter, r *http.Request) {
	var req VerifyCodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	fullMobile := req.Country + req.Mobile
	attemptKey := "verify_limit:" + fullMobile

	// Check if attempts exceeded
	attempts, _ := redisClient.Get(ctx, attemptKey).Int()
	if attempts >= 3 {
		http.Error(w, "too many wrong attempts, request new code", http.StatusTooManyRequests)
		return
	}

	storedCode, err := redisClient.Get(ctx, "sms:"+fullMobile).Result()
	if err == redis.Nil || storedCode != req.Code {
		// Increment failed attempts
		redisClient.Incr(ctx, attemptKey)
		redisClient.Expire(ctx, attemptKey, 2*time.Minute)
		http.Error(w, "invalid or expired code", http.StatusUnauthorized)
		return
	}

	// Successful verify: reset attempt counter and delete SMS code
	redisClient.Del(ctx, attemptKey)
	redisClient.Del(ctx, "sms:"+fullMobile)

	// Check if user exists
	userColl := mongoClient.Database(os.Getenv("DB_NAME")).Collection("users")
	var user models.User
	newUser := false
	if err := userColl.FindOne(ctx, bson.M{"mobile": fullMobile}).Decode(&user); err == mongo.ErrNoDocuments {
		newUser = true
	} else if err != nil {
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}

	// Generate Access Token
	accessClaims := jwt.MapClaims{"mobile": fullMobile, "exp": time.Now().Add(15 * time.Minute).Unix()}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	signedAccess, err := accessToken.SignedString(jwtSecret)
	if err != nil {
		http.Error(w, "cannot generate access token", http.StatusInternalServerError)
		return
	}
	redisClient.Set(ctx, "token:"+fullMobile, signedAccess, 15*time.Minute)

	// Generate Refresh Token
	refreshClaims := jwt.MapClaims{"mobile": fullMobile, "exp": time.Now().Add(7 * 24 * time.Hour).Unix()}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	signedRefresh, err := refreshToken.SignedString(jwtSecret)
	if err != nil {
		http.Error(w, "cannot generate refresh token", http.StatusInternalServerError)
		return
	}
	redisClient.Set(ctx, "refresh:"+fullMobile, signedRefresh, 7*24*time.Hour)

	// Return tokens
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"accessToken":  signedAccess,
		"refreshToken": signedRefresh,
		"newUser":      newUser,
	})
}

// Refresh Access Token using Refresh Token
func RefreshToken(w http.ResponseWriter, r *http.Request) {
	var body struct {
		RefreshToken string `json:"refreshToken"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	token, err := jwt.Parse(body.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		http.Error(w, "invalid refresh token", http.StatusUnauthorized)
		return
	}

	claims := token.Claims.(jwt.MapClaims)
	mobile, ok := claims["mobile"].(string)
	if !ok || mobile == "" {
		http.Error(w, "invalid refresh token", http.StatusUnauthorized)
		return
	}

	val, err := redisClient.Get(ctx, "refresh:"+mobile).Result()
	if err != nil || val != body.RefreshToken {
		http.Error(w, "refresh token expired", http.StatusUnauthorized)
		return
	}

	newClaims := jwt.MapClaims{"mobile": mobile, "exp": time.Now().Add(15 * time.Minute).Unix()}
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)
	signed, err := newToken.SignedString(jwtSecret)
	if err != nil {
		http.Error(w, "cannot generate new access token", http.StatusInternalServerError)
		return
	}

	redisClient.Set(ctx, "token:"+mobile, signed, 15*time.Minute)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"accessToken": signed})
}

// Validate Access Token
func ValidateToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		json.NewEncoder(w).Encode(map[string]any{"loggedIn": false, "expired": false})
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) { return jwtSecret, nil })
	if err != nil {
		expired := errors.Is(err, jwt.ErrTokenExpired)
		json.NewEncoder(w).Encode(map[string]any{"loggedIn": false, "expired": expired})
		return
	}
	if !token.Valid {
		json.NewEncoder(w).Encode(map[string]any{"loggedIn": false, "expired": false})
		return
	}

	claims := token.Claims.(jwt.MapClaims)
	mobile, ok := claims["mobile"].(string)
	if !ok || mobile == "" {
		json.NewEncoder(w).Encode(map[string]any{"loggedIn": false, "expired": false})
		return
	}

	val, err := redisClient.Get(ctx, "token:"+mobile).Result()
	if err != nil || val != tokenString {
		json.NewEncoder(w).Encode(map[string]any{"loggedIn": false, "expired": false})
		return
	}

	userColl := mongoClient.Database(os.Getenv("DB_NAME")).Collection("users")
	var user models.User
	userHasInfo := true
	if err := userColl.FindOne(ctx, bson.M{"mobile": mobile}).Decode(&user); err == mongo.ErrNoDocuments {
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

// Logout user
func Logout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		json.NewEncoder(w).Encode(map[string]any{"loggedIn": false, "Logout": true})
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) { return jwtSecret, nil })
	if token != nil && token.Valid {
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if mobile, ok := claims["mobile"].(string); ok && mobile != "" {
				redisClient.Del(ctx, "token:"+mobile)
				redisClient.Del(ctx, "refresh:"+mobile)
			}
		}
	}

	json.NewEncoder(w).Encode(map[string]any{"loggedIn": false, "Logout": true})
}

func Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		http.Error(w, `{"success":false,"error":"missing token"}`, http.StatusUnauthorized)
		return
	}
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]any{"success": false, "error": "invalid token"})
		return
	}
	claims := token.Claims.(jwt.MapClaims)
	mobile, ok := claims["mobile"].(string)
	if !ok || mobile == "" {
		http.Error(w, `{"success":false,"error":"invalid token claims"}`, http.StatusUnauthorized)
		return
	}
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, `{"success":false,"error":"cannot parse form"}`, http.StatusBadRequest)
		return
	}
	name := r.FormValue("name")
	family := r.FormValue("family")
	if name == "" || family == "" {
		http.Error(w, `{"success":false,"error":"name/family required"}`, http.StatusBadRequest)
		return
	}
	var photoPath string
	file, header, err := r.FormFile("photo")
	if err == nil && file != nil {
		defer file.Close()
		fileName := fmt.Sprintf("%s_%d_%s", mobile, time.Now().Unix(), header.Filename)
		url, err := service.MinioUpload(fileName, file, header.Size, header.Header.Get("Content-Type"))
		if err != nil {
			http.Error(w, `{"success":false,"error":"cannot upload photo"}`, http.StatusInternalServerError)
			return
		}
		photoPath = url
		fmt.Println("File removed from MinIO:", url)
	}
	userColl := mongoClient.Database(os.Getenv("DB_NAME")).Collection("users")
	newUser := models.User{
		Name:         name,
		Family:       family,
		Mobile:       mobile,
		Balance:      0,
		TokenBalance: 0,
		Image:        photoPath,
		FingerTokens: []models.WebAuthnCredential{},
	}
	_, err = userColl.InsertOne(ctx, newUser)
	if err != nil {
		http.Error(w, `{"success":false,"error":"database error"}`, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"success":true}`))
}
