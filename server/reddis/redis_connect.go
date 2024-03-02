package reddis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client
var PubSubConnection *redis.PubSub

// type RedisConfig struct{

// }

// func init(){

// }

// TODO TO ADD CONFIG FILE

func InitRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	ctx := context.Background()
	PubSubConnection = RedisClient.Subscribe(ctx)

	if err := RedisClient.Ping(ctx).Err(); err != nil {
		panic(err)
	}
}
