package db

import (
	"context"
	"fmt"
	"log"

	"github.com/SubhamMurarka/chat_app/server/config"
	"github.com/redis/go-redis/v9"
)

type RedisDatabase struct {
	client *redis.Client
	pubSub *redis.PubSub
}

func NewRedisDatabase() (*RedisDatabase, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.Config.RedisHost, config.Config.RedisPort),
		Password: "",
		DB:       0,
		PoolSize: 5,
	})

	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	err := client.ConfigSet(ctx, "notify-keyspace-events", "Ex").Err()
	if err != nil {
		log.Println("Error configuring for notifying ttl")
		return nil, err
	}

	pubSub := client.Subscribe(ctx)
	pubSub.Subscribe(ctx, "__keyevent@0__:expired")

	return &RedisDatabase{client: client, pubSub: pubSub}, nil
}

func (r *RedisDatabase) GetClient() *redis.Client {
	return r.client
}

func (r *RedisDatabase) GetPubSub() *redis.PubSub {
	return r.pubSub
}

func (r *RedisDatabase) Close() {
	r.client.Close()
	r.pubSub.Close()
}
