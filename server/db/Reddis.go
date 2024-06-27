package db

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type RedisDatabase struct {
	client *redis.Client
	pubSub *redis.PubSub
}

func NewRedisDatabase() (*RedisDatabase, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	pubSub := client.Subscribe(ctx)

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
