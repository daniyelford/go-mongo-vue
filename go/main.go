package main

import (
	users "go-mongo-vue-users"
	"log"
	"net/http"
	"os"
)

func main() {
	r, err := users.Init()
	if err != nil {
		log.Fatal("error:", err)
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server running on :%s\n", port)
	http.ListenAndServe(":"+port, r)
}
