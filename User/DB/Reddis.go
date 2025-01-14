package db

import (
	"context"
	"fmt"

	"github.com/SubhamMurarka/chat_app/User/Config"
	"github.com/redis/go-redis/v9"
)

type RedisDatabase struct {
	client *redis.Client
	pubSub *redis.PubSub
}

func NewRedisDatabase() (*RedisDatabase, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", Config.Config.RedisHost, Config.Config.RedisPort),
		Password: "",
		DB:       0,
		PoolSize: 50,
	})

	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	pubSub := client.Subscribe(ctx)

	return &RedisDatabase{pubSub: pubSub}, nil
}

func (r *RedisDatabase) GetPubSub() *redis.PubSub {
	return r.pubSub
}

func (r *RedisDatabase) GetClient() *redis.Client {
	return r.client
}

func (r *RedisDatabase) Close() {
	r.pubSub.Close()
	r.client.Close()
}
