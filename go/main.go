package main

import (
	"context"
	"fmt"
	"go-mongo-vue-go/handlers"
	"go-mongo-vue-go/router"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ctx = context.Background()

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("⚠️ no .env file found, using system env")
	}
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("Mongo ping failed:", err)
	}
	handlers.InitSMSCollection(client, os.Getenv("DB_NAME"))
	r := router.NewRouter()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server running on :%s\n", port)
	http.ListenAndServe(":"+port, r)
}
