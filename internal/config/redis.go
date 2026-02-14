package config

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client

var ctx = context.Background()

func ConnectRedis() {
	redisURL := strings.TrimSpace(os.Getenv("REDIS_URL"))

	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		panic(fmt.Sprintf("Invalid REDIS_URL: %v", err))
	}
	RDB = redis.NewClient(opt)
	// RDB = redis.NewClient(&redis.Options{
	// 	Addr: "127.0.0.1:6379", // use IPv4 instead of [::1]
	// 	// Addr: "localhost:6379", // Ideally, load this from env variables
	// 	// Addr:     "redis:6379", // redis is the service name in docker-compose
	// 	Password: "",
	// 	DB:       0,
	// })

	// Test the connection immediately
	if err := RDB.Ping(ctx).Err(); err != nil {
		panic(fmt.Sprintf("Failed to connect to Redis: %v", err))
	}

	fmt.Println("Connected to Redis successfully")
}
