package users

import (
	"context"
	"fmt"
	"go-mongo-vue-users/config"
	"go-mongo-vue-users/router"
	"go-mongo-vue-users/service"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ctx = context.Background()

func Init() (*mux.Router, error) {
	if err := godotenv.Load(); err != nil {
		fmt.Println("⚠️ no .env file found, using system env")
	}
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil {
		return nil, err
		// log.Fatal(err)
	}
	if err := client.Ping(ctx, nil); err != nil {
		// log.Fatal("Mongo ping failed:", err)
		return nil, err
	}
	config.InitMongo(client)
	redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASS"),
		DB:       0,
	})
	if err := redisClient.Ping(ctx).Err(); err != nil {
		// log.Fatal("Redis ping failed:", err)
		return nil, err
	}
	config.InitRedis(redisClient)
	if err := service.MinioInit(); err != nil {
		// log.Fatal("MinIO init error:", err)
		return nil, err
	}
	service.WebAuthnInit()
	return router.NewRouter(), nil
}
