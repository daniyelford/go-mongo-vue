package handlers

import (
	"encoding/json"
	"go-mongo-vue-go/models"
	"go-mongo-vue-go/service"
	"net/http"

	"github.com/go-webauthn/webauthn/webauthn"
)

var sessions = make(map[string]*webauthn.SessionData)

func LoginFingerPrintStart(w http.ResponseWriter, r *http.Request) {
	// فرض: یوزر رو از DB با شماره موبایل پیدا می‌کنیم
	user := models.User{
		Mobile:       "0912000000",
		FingerTokens: []models.WebAuthnCredential{
			// اینجا credential های قبلی رو از DB load کن
		},
	}
	options, session, err := service.WA.BeginLogin(&user)
	if err != nil {
		http.Error(w, "failed begin login", http.StatusInternalServerError)
		return
	}
	sessions[user.Mobile] = session
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(options)
}
func LoginFingerPrintEnd(w http.ResponseWriter, r *http.Request) {
	user := models.User{
		Mobile:       "0912000000", // در عمل از DB load کن
		FingerTokens: []models.WebAuthnCredential{
			// credential های قبلی
		},
	}
	session := sessions[user.Mobile]
	if session == nil {
		http.Error(w, "session not found", http.StatusBadRequest)
		return
	}
	_, err := service.WA.FinishLogin(&user, *session, r)
	if err != nil {
		http.Error(w, "login failed: "+err.Error(), http.StatusUnauthorized)
		return
	}
	delete(sessions, user.Mobile)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"message": "login successful",
	})
}
func RegisterFingerPrintStart(w http.ResponseWriter, r *http.Request) {
	// فرض: یوزر رو از DB لود می‌کنی
	user := models.User{
		Name:   "Ali",
		Family: "Rezayi",
		Mobile: "0912000000",
	}
	options, session, err := service.WA.BeginRegistration(&user)
	if err != nil {
		http.Error(w, "failed begin registration", http.StatusInternalServerError)
		return
	}
	sessions[user.Mobile] = session
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(options)
}
func RegisterFingerPrintEnd(w http.ResponseWriter, r *http.Request) {
	user := models.User{
		Mobile: "0912000000", // در عمل باید از DB لود بشه
	}
	session := sessions[user.Mobile]
	if session == nil {
		http.Error(w, "session not found", http.StatusBadRequest)
		return
	}
	credential, err := service.WA.FinishRegistration(&user, *session, r)
	if err != nil {
		http.Error(w, "finish registration failed: "+err.Error(), http.StatusBadRequest)
		return
	}
	cred := models.WebAuthnCredential{
		CredentialID: credential.ID,
		PublicKey:    credential.PublicKey,
		SignCount:    credential.Authenticator.SignCount,
	}
	user.FingerTokens = append(user.FingerTokens, cred)

	// TODO: ذخیره user در MongoDB
	// db.Users.UpdateOne(...)

	// پاک کردن session بعد از استفاده
	delete(sessions, user.Mobile)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"message": "registered successfully",
	})
}
