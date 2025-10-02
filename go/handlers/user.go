package handlers

import (
	"encoding/json"
	"fmt"
	"go-mongo-vue-go/config"
	"go-mongo-vue-go/libraries"
	"go-mongo-vue-go/middleware"
	"go-mongo-vue-go/models"
	"go-mongo-vue-go/service"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type SendCodeRequest struct {
	Mobile  string `json:"mobile"`
	Country string `json:"country"`
}
type VerifyCodeRequest struct {
	Mobile  string `json:"mobile"`
	Country string `json:"country"`
	Code    string `json:"code"`
}

func SendCode(w http.ResponseWriter, r *http.Request) {
	var req SendCodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	fullMobile := req.Country + req.Mobile
	limitKey := "sms_limit:" + fullMobile
	lastKey := "sms_last:" + fullMobile
	lastTs, _ := config.RedisClient.Get(config.Ctx, lastKey).Int64()
	if time.Now().Unix()-lastTs < 15 {
		http.Error(w, "please wait 15 seconds before retrying", http.StatusTooManyRequests)
		return
	}
	count, _ := config.RedisClient.Get(config.Ctx, limitKey).Int()
	if count >= 3 {
		http.Error(w, "too many requests, try again in 5 minutes", http.StatusTooManyRequests)
		return
	}
	code := fmt.Sprintf("%06d", rand.Intn(900000)+100000)
	if err := config.RedisClient.Set(config.Ctx, "sms:"+fullMobile, code, 5*time.Minute).Err(); err != nil {
		http.Error(w, "failed to save code", http.StatusInternalServerError)
		return
	}
	if os.Getenv("SMS_SANDBOX") != "true" {
		ok, err := libraries.SendSmsLogin(code, fullMobile)
		if err != nil || !ok {
			http.Error(w, "failed to send sms", http.StatusInternalServerError)
			return
		}
	}
	config.RedisClient.Incr(config.Ctx, limitKey)
	config.RedisClient.Expire(config.Ctx, limitKey, 5*time.Minute)
	config.RedisClient.Set(config.Ctx, lastKey, time.Now().Unix(), 5*time.Minute)
	result := map[string]string{"message": "code sent"}
	if os.Getenv("SMS_SANDBOX") == "true" {
		result["code"] = code
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
func VerifyCode(w http.ResponseWriter, r *http.Request) {
	var req VerifyCodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	fullMobile := req.Country + req.Mobile
	attemptKey := "verify_limit:" + fullMobile
	attempts, _ := config.RedisClient.Get(config.Ctx, attemptKey).Int()
	if attempts >= 3 {
		http.Error(w, "too many wrong attempts, request new code", http.StatusTooManyRequests)
		return
	}
	storedCode, err := config.RedisClient.Get(config.Ctx, "sms:"+fullMobile).Result()
	if err == redis.Nil || storedCode != req.Code {
		config.RedisClient.Incr(config.Ctx, attemptKey)
		config.RedisClient.Expire(config.Ctx, attemptKey, 2*time.Minute)
		http.Error(w, "invalid or expired code", http.StatusUnauthorized)
		return
	}
	config.RedisClient.Del(config.Ctx, attemptKey)
	config.RedisClient.Del(config.Ctx, "sms:"+fullMobile)
	userColl := config.MongoClient.Database(os.Getenv("DB_NAME")).Collection("users")
	var user models.User
	newUser := false
	if err := userColl.FindOne(config.Ctx, bson.M{"mobile": fullMobile}).Decode(&user); err == mongo.ErrNoDocuments {
		newUser = true
	} else if err != nil {
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}
	accessClaims := jwt.MapClaims{"mobile": fullMobile, "exp": time.Now().Add(15 * time.Minute).Unix()}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	signedAccess, err := accessToken.SignedString(config.JwtSecret)
	if err != nil {
		http.Error(w, "cannot generate access token", http.StatusInternalServerError)
		return
	}
	config.RedisClient.Set(config.Ctx, "token:"+fullMobile, signedAccess, 15*time.Minute)
	refreshClaims := jwt.MapClaims{"mobile": fullMobile, "exp": time.Now().Add(7 * 24 * time.Hour).Unix()}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	signedRefresh, err := refreshToken.SignedString(config.JwtSecret)
	if err != nil {
		http.Error(w, "cannot generate refresh token", http.StatusInternalServerError)
		return
	}
	config.RedisClient.Set(config.Ctx, "refresh:"+fullMobile, signedRefresh, 7*24*time.Hour)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"accessToken":  signedAccess,
		"refreshToken": signedRefresh,
		"newUser":      newUser,
	})
}
func Logout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	mobile := r.Context().Value(middleware.MobileContextKey).(string)
	config.RedisClient.Del(config.Ctx, "token:"+mobile)
	config.RedisClient.Del(config.Ctx, "refresh:"+mobile)
	json.NewEncoder(w).Encode(map[string]any{"loggedIn": false, "Logout": true})
}
func Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	mobile := r.Context().Value(middleware.MobileContextKey).(string)
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
	var fileName string
	file, header, err := r.FormFile("photo")
	if err == nil && file != nil {
		defer file.Close()
		fileName = fmt.Sprintf("%s_%d_%s", mobile, time.Now().Unix(), header.Filename)
		if _, err := service.MinioUpload(fileName, file, header.Size, header.Header.Get("Content-Type")); err != nil {
			http.Error(w, `{"success":false,"error":"cannot upload photo"}`, http.StatusInternalServerError)
			return
		}
	}
	userColl := config.MongoClient.Database(os.Getenv("DB_NAME")).Collection("users")
	newUser := models.User{
		Name:         name,
		Family:       family,
		Mobile:       mobile,
		Balance:      0,
		TokenBalance: 0,
		Image:        fileName,
		FingerTokens: []models.WebAuthnCredential{},
	}
	_, err = userColl.InsertOne(config.Ctx, newUser)
	if err != nil {
		http.Error(w, `{"success":false,"error":"database error"}`, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"success":true}`))
}
func UserInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	mobile := r.Context().Value(middleware.MobileContextKey).(string)
	userColl := config.MongoClient.Database(os.Getenv("DB_NAME")).Collection("users")
	var user models.User
	if err := userColl.FindOne(config.Ctx, bson.M{"mobile": mobile}).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, `{"success":false,"error":"user not found"}`, http.StatusNotFound)
		} else {
			http.Error(w, `{"success":false,"error":"database error"}`, http.StatusInternalServerError)
		}
		return
	}
	var imageURL string
	if user.Image != "" {
		publicEndpoint := os.Getenv("MINIO_PUBLIC_ENDPOINT")
		if publicEndpoint == "" {
			publicEndpoint = os.Getenv("MINIO_ENDPOINT")
		}
		imageURL = fmt.Sprintf(
			"http://%s/%s/%s",
			publicEndpoint,
			os.Getenv("MINIO_BUCKET"),
			url.PathEscape(user.Image),
		)
	}
	json.NewEncoder(w).Encode(map[string]any{
		"success": true,
		"user": map[string]any{
			"name":         user.Name,
			"family":       user.Family,
			"mobile":       user.Mobile,
			"balance":      user.Balance,
			"tokenBalance": user.TokenBalance,
			"image":        imageURL,
			"fingerTokens": user.FingerTokens,
		},
	})
}
func UserUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	mobile := r.Context().Value(middleware.MobileContextKey).(string)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, `{"success":false,"error":"cannot parse form"}`, http.StatusBadRequest)
		return
	}
	name := r.FormValue("name")
	family := r.FormValue("family")
	updateFields := bson.M{}
	if name != "" {
		updateFields["name"] = name
	}
	if family != "" {
		updateFields["family"] = family
	}
	file, header, err := r.FormFile("photo")
	if err == nil && file != nil {
		defer file.Close()
		fileName := fmt.Sprintf("%s_%d_%s", mobile, time.Now().Unix(), header.Filename)
		if _, err := service.MinioUpload(fileName, file, header.Size, header.Header.Get("Content-Type")); err != nil {
			http.Error(w, `{"success":false,"error":"cannot upload photo"}`, http.StatusInternalServerError)
			return
		}
		updateFields["image"] = fileName
	}

	if len(updateFields) == 0 {
		http.Error(w, `{"success":false,"error":"no fields to update"}`, http.StatusBadRequest)
		return
	}

	userColl := config.MongoClient.Database(os.Getenv("DB_NAME")).Collection("users")
	_, err = userColl.UpdateOne(
		config.Ctx,
		bson.M{"mobile": mobile},
		bson.M{"$set": updateFields},
	)
	if err != nil {
		http.Error(w, `{"success":false,"error":"database error"}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]any{
		"success": true,
		"updated": updateFields,
	})
}
