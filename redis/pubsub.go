package redis

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	RedisClient *redis.Client
	Ctx         = context.Background()
)

func InitRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})
	_, err := RedisClient.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
}

func PublishMessage(channel, message string) error {
	ctx, cancel := context.WithTimeout(Ctx, 5*time.Second)
	defer cancel()
	return RedisClient.Publish(ctx, channel, message).Err()
}

func Subscribe(channel string) *redis.PubSub {
	return RedisClient.Subscribe(Ctx, channel)
}
