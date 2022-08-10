package config

import (
	"log"

	redis "github.com/go-redis/redis/v8"
)

var Redis *redis.Client

// Redis client creation functionality
func CreateRedisClient() {
	opt, err := redis.ParseURL("redis://localhost:6379/0")
	if err != nil {
		log.Fatalf("Redis Connection Failed to Open: %v", err.Error())
	}

	Redis = redis.NewClient(opt)

}
