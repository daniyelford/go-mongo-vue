package handlers

import (
	"encoding/json"
	"fmt"
	sms "go-mongo-vue-back/libraries"
	"math/rand"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("super-secret-key")
var codes = map[string]string{}

type SendCodeRequest struct {
	Mobile string `json:"mobile"`
}
type VerifyCodeRequest struct {
	Mobile string `json:"mobile"`
	Code   string `json:"code"`
}
type LoginResponse struct {
	Token string `json:"token"`
}

func SendCode(w http.ResponseWriter, r *http.Request) {
	var req SendCodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	code := rand.Intn(900000) + 100000
	codes[req.Mobile] = fmt.Sprintf("%d", code)
	ok, err := sms.SendSmsForce(codes[req.Mobile], req.Mobile)
	if err != nil || !ok {
		http.Error(w, "failed to send sms", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "code sent",
		// "code": codes[req.Mobile], // ❌ فقط برای تست می‌تونی موقت فعال کنی
	})
}
func VerifyCode(w http.ResponseWriter, r *http.Request) {
	var req VerifyCodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	if codes[req.Mobile] != req.Code {
		http.Error(w, "invalid code", http.StatusUnauthorized)
		return
	}
	claims := jwt.MapClaims{
		"mobile": req.Mobile,
		"exp":    time.Now().Add(time.Minute * 15).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(jwtSecret)
	if err != nil {
		http.Error(w, "could not generate token", http.StatusInternalServerError)
		return
	}
	delete(codes, req.Mobile)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(LoginResponse{Token: signed})
}
