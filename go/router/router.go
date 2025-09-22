package router

import (
	"go-mongo-vue-go/handlers"

	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/api/login/send", handlers.SendCode).Methods("POST")
	r.HandleFunc("/api/login/verify", handlers.VerifyCode).Methods("POST")

	r.HandleFunc("/api/login/fingerPrint/start", handlers.LoginFingerPrintStart).Methods("POST")
	r.HandleFunc("/api/login/fingerPrint/end", handlers.LoginFingerPrintEnd).Methods("POST")
	r.HandleFunc("/api/register/fingerPrint/start", handlers.RegisterFingerPrintStart).Methods("POST")
	r.HandleFunc("/api/register/fingerPrint/end", handlers.RegisterFingerPrintEnd).Methods("POST")

	r.HandleFunc("/api/auth/logout", handlers.Logout).Methods("GET")

	r.HandleFunc("/api/auth/validate", handlers.ValidateToken).Methods("GET")
	r.HandleFunc("/api/token/refresh", handlers.RefreshToken).Methods("GET")

	r.HandleFunc("/api/hello", handlers.HelloWorld).Methods("GET")
	return r
}
