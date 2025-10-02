package router

import (
	"go-mongo-vue-go/handlers"
	"go-mongo-vue-go/middleware"

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
	r.HandleFunc("/api/token/refresh", handlers.RefreshToken).Methods("GET")
	private := r.PathPrefix("/api").Subrouter()
	private.Use(middleware.JWTAuthMiddleware)
	private.HandleFunc("/auth/logout", handlers.Logout).Methods("GET")
	private.HandleFunc("/auth/validate", handlers.ValidateToken).Methods("GET")
	private.HandleFunc("/register/save", handlers.Register).Methods("POST")
	private.HandleFunc("/user/update", handlers.UserUpdate).Methods("POST")
	private.HandleFunc("/user/info", handlers.UserInfo).Methods("GET")
	return r
}
