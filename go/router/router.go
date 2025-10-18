package router

import (
	"go-mongo-vue-go/config"
	"go-mongo-vue-go/handlers"
	"go-mongo-vue-go/middleware"

	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/health", config.Health).Methods("GET", "HEAD")
	r.HandleFunc("/api/login/fingerPrint/has", handlers.HasFingerPrint).Methods("POST")
	r.HandleFunc("/api/login/send", handlers.SendCode).Methods("POST")
	r.HandleFunc("/api/login/verify", handlers.VerifyCode).Methods("POST")
	r.HandleFunc("/api/login/fingerPrint/start", handlers.LoginFingerPrintStart).Methods("POST")
	r.HandleFunc("/api/login/fingerPrint/end", handlers.LoginFingerPrintEnd).Methods("POST")
	r.HandleFunc("/api/token/refresh", handlers.RefreshToken).Methods("POST")
	private := r.PathPrefix("/api").Subrouter()
	private.Use(middleware.JWTAuthMiddleware)
	private.HandleFunc("/register/fingerPrint/start", handlers.RegisterFingerPrintStart).Methods("POST")
	private.HandleFunc("/register/fingerPrint/end", handlers.RegisterFingerPrintEnd).Methods("POST")
	private.HandleFunc("/auth/logout", handlers.Logout).Methods("GET")
	private.HandleFunc("/auth/validate", handlers.ValidateToken).Methods("GET")
	private.HandleFunc("/register/save", handlers.Register).Methods("POST")
	private.HandleFunc("/user/update", handlers.UserUpdate).Methods("PUT")
	private.HandleFunc("/user/info", handlers.UserInfo).Methods("GET")
	private.HandleFunc("/posts/all", handlers.GetAllPosts).Methods("POST")
	private.HandleFunc("/posts/create", handlers.CreatePost).Methods("POST")
	private.HandleFunc("/posts/edit", handlers.EditPost).Methods("PUT")
	private.HandleFunc("/posts/delete", handlers.DeletePost).Methods("DELETE")
	return r
}
