package router

import (
	"go-vue-nosql-back/handlers"

	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/api/users", handlers.GetUsers).Methods("GET")
	r.HandleFunc("/api/hello", handlers.HelloWorld).Methods("GET")
	return r
}
