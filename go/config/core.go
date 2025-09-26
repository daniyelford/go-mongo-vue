package config

import (
	"context"
	"os"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	JwtSecret   = []byte(os.Getenv("JWT_SECRET"))
	RedisClient *redis.Client
	Ctx         = context.Background()
	MongoClient *mongo.Client
)

func InitRedis(client *redis.Client) { RedisClient = client }
func InitMongo(client *mongo.Client) { MongoClient = client }
