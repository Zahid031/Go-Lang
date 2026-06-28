package redisclient

import (
	"context"
	"log"
	"github.com/redis/go-redis/v9"
	"os"
)

func NewClient(ctx context.Context) *redis.Client {
	redisAddr := os.Getenv("REDIS_ADDR")
	redisDB := 0 // default DB

	client := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		DB:       redisDB,
	})

	// Test the connection
	if err := client.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	return client
}