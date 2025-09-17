package main

import (
	"fmt"
	"go-mongo-vue-go/router"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("⚠️ no .env file found, using system env")
	}
	r := router.NewRouter()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server running on :%s\n", port)
	http.ListenAndServe(":"+port, r)
}
