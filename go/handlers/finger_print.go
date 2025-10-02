package handlers

import (
	"encoding/json"
	"go-mongo-vue-go/config"
	"go-mongo-vue-go/middleware"
	"go-mongo-vue-go/models"
	"go-mongo-vue-go/service"
	"net/http"
	"os"
	"time"

	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func saveSessionToRedis(mobile string, session *webauthn.SessionData) error {
	data, err := json.Marshal(session)
	if err != nil {
		return err
	}
	return config.RedisClient.Set(config.Ctx, "webauthn_session:"+mobile, data, 5*time.Minute).Err()
}
func loadSessionFromRedis(mobile string) (*webauthn.SessionData, error) {
	data, err := config.RedisClient.Get(config.Ctx, "webauthn_session:"+mobile).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var session webauthn.SessionData
	if err := json.Unmarshal(data, &session); err != nil {
		return nil, err
	}
	return &session, nil
}
func deleteSessionFromRedis(mobile string) {
	config.RedisClient.Del(config.Ctx, "webauthn_session:"+mobile)
}
func findUserByMobile(mobile string) (*models.User, error) {
	userColl := config.MongoClient.Database(os.Getenv("DB_NAME")).Collection("users")
	var user models.User
	if err := userColl.FindOne(config.Ctx, bson.M{"mobile": mobile}).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
func RegisterFingerPrintStart(w http.ResponseWriter, r *http.Request) {
	mobile := r.Context().Value(middleware.MobileContextKey).(string)
	user, err := findUserByMobile(mobile)
	if err != nil || user == nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}
	options, session, err := service.WA.BeginRegistration(user)
	if err != nil {
		http.Error(w, "failed begin registration", http.StatusInternalServerError)
		return
	}
	if err := saveSessionToRedis(user.Mobile, session); err != nil {
		http.Error(w, "failed to save session", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(options)
}
func RegisterFingerPrintEnd(w http.ResponseWriter, r *http.Request) {
	mobile := r.Context().Value(middleware.MobileContextKey).(string)
	user, err := findUserByMobile(mobile)
	if err != nil || user == nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}
	session, err := loadSessionFromRedis(user.Mobile)
	if err != nil || session == nil {
		http.Error(w, "session not found", http.StatusBadRequest)
		return
	}
	credential, err := service.WA.FinishRegistration(user, *session, r)
	if err != nil {
		http.Error(w, "finish registration failed: "+err.Error(), http.StatusBadRequest)
		return
	}
	cred := models.WebAuthnCredential{
		CredentialID: credential.ID,
		PublicKey:    credential.PublicKey,
		SignCount:    credential.Authenticator.SignCount,
	}
	userColl := config.MongoClient.Database(os.Getenv("DB_NAME")).Collection("users")
	_, err = userColl.UpdateOne(
		config.Ctx,
		bson.M{"mobile": user.Mobile},
		bson.M{"$push": bson.M{"finger_tokens": cred}},
	)
	if err != nil {
		http.Error(w, "failed to save credential", http.StatusInternalServerError)
		return
	}
	deleteSessionFromRedis(user.Mobile)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"message": "registered successfully",
	})
}
func LoginFingerPrintStart(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Mobile string `json:"mobile"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Mobile == "" {
		http.Error(w, "mobile required", http.StatusBadRequest)
		return
	}
	if body.Mobile == "" {
		http.Error(w, "mobile required", http.StatusBadRequest)
		return
	}
	user, err := findUserByMobile(body.Mobile)
	if err != nil || user == nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}
	options, session, err := service.WA.BeginLogin(user)
	if err != nil {
		http.Error(w, "failed begin login", http.StatusInternalServerError)
		return
	}
	if err := saveSessionToRedis(user.Mobile, session); err != nil {
		http.Error(w, "failed to save session", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(options)
}
func LoginFingerPrintEnd(w http.ResponseWriter, r *http.Request) {
	mobile := r.Header.Get("X-Mobile")
	if mobile == "" {
		http.Error(w, "mobile header required", http.StatusBadRequest)
		return
	}
	user, err := findUserByMobile(mobile)
	if err != nil || user == nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}
	session, err := loadSessionFromRedis(user.Mobile)
	if err != nil || session == nil {
		http.Error(w, "session not found", http.StatusBadRequest)
		return
	}
	credential, err := service.WA.FinishLogin(user, *session, r)
	if err != nil {
		http.Error(w, "login failed: "+err.Error(), http.StatusUnauthorized)
		return
	}
	userColl := config.MongoClient.Database(os.Getenv("DB_NAME")).Collection("users")
	_, err = userColl.UpdateOne(
		config.Ctx,
		bson.M{"mobile": user.Mobile, "finger_tokens.credential_id": credential.ID},
		bson.M{"$set": bson.M{"finger_tokens.$.sign_count": credential.Authenticator.SignCount}},
	)
	if err != nil {
		http.Error(w, "failed to update sign_count", http.StatusInternalServerError)
		return
	}
	deleteSessionFromRedis(user.Mobile)
	accessClaims := jwt.MapClaims{"mobile": user.Mobile, "exp": time.Now().Add(15 * time.Minute).Unix()}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	signedAccess, err := accessToken.SignedString(config.JwtSecret)
	if err != nil {
		http.Error(w, "cannot generate access token", http.StatusInternalServerError)
		return
	}
	config.RedisClient.Set(config.Ctx, "token:"+user.Mobile, signedAccess, 15*time.Minute)
	refreshClaims := jwt.MapClaims{"mobile": user.Mobile, "exp": time.Now().Add(7 * 24 * time.Hour).Unix()}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	signedRefresh, err := refreshToken.SignedString(config.JwtSecret)
	if err != nil {
		http.Error(w, "cannot generate refresh token", http.StatusInternalServerError)
		return
	}
	config.RedisClient.Set(config.Ctx, "refresh:"+user.Mobile, signedRefresh, 7*24*time.Hour)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"accessToken":  signedAccess,
		"refreshToken": signedRefresh,
	})
}
func HasFingerPrint(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Mobile string `json:"mobile"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Mobile == "" {
		http.Error(w, "mobile required", http.StatusBadRequest)
		return
	}
	if body.Mobile == "" {
		http.Error(w, "mobile required", http.StatusBadRequest)
		return
	}
	user, err := findUserByMobile(body.Mobile)
	if err != nil {
		http.Error(w, "error finding user", http.StatusInternalServerError)
		return
	}
	if user == nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}
	status := false
	if len(user.FingerTokens) > 0 {
		status = true
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"status": status,
	})
}
