package redis

import (
	"context"
	"fmt"
	"log"

	"github.com/SubhamMurarka/chat_app/Image/Config"
	"github.com/redis/go-redis/v9"
)

var Client *redis.Client

func NewRedisDatabase() error {
	Client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", Config.Conf.RedisHost, Config.Conf.RedisPort),
		Password: "",
		DB:       0,
		PoolSize: 5,
	})

	ctx := context.Background()
	if err := Client.Ping(ctx).Err(); err != nil {
		log.Fatalf("error connecting redis")
		return err
	}

	log.Printf("Redis Connected!")

	return nil
}

func IsUserActive(key string) bool {
	_, err := Client.Get(context.Background(), key).Result()

	if err == redis.Nil {
		fmt.Println("User is not online")
		return false
	} else if err != nil {
		log.Fatalf("Could not fetch data from Redis: %v", err)
		return false
	}
	return true
}
