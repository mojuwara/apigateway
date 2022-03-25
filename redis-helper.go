package main

import (
	"context"

	"github.com/go-redis/redis/v8"
)

// TODO: Update to read fields from files set in env variables
func initRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	// Check if we can ping Redis successfully
	ctx := context.Background()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		logger.Fatal("Error while connecting to Redis:", err)
	} else {
		logger.Println("Successfully connected to Redis")
	}
}
