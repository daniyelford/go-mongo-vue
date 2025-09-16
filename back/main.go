package main

import (
	"go-mongo-vue-back/router"
	"log"
	"net/http"
)

func main() {
	r := router.NewRouter()
	log.Println("Server running on :8080")
	http.ListenAndServe(":8080", r)
}
