package router

import (
	"go-mongo-vue-back/handlers"

	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/api/login/send", handlers.SendCode).Methods("POST")
	r.HandleFunc("/api/login/verify", handlers.VerifyCode).Methods("POST")

	r.HandleFunc("/api/hello", handlers.HelloWorld).Methods("GET")
	return r
}
