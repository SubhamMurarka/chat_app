package Repository

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/SubhamMurarka/chat_app/Heartbeat/Models"
	"github.com/redis/go-redis/v9"
)

type HeartBeatRepository interface {
	UserTTlManager(ctx context.Context, clientID int64, channelID int64) error
	ListenMessages()
}

type heartbeatRepository struct {
	redisClient *redis.Client
	redisPubSub *redis.PubSub
}

func NewHeartBeatRepository(redisClient *redis.Client, redisPubSub *redis.PubSub) HeartBeatRepository {
	return &heartbeatRepository{redisClient: redisClient, redisPubSub: redisPubSub}
}

func (r *heartbeatRepository) UserTTlManager(ctx context.Context, clientID int64, channelID int64) error {
	ttl := 50 * time.Second

	cl := strconv.FormatInt(clientID, 10)
	ch := strconv.FormatInt(channelID, 10)

	fmt.Printf("client : %s, channel : %s", cl, ch)

	key := fmt.Sprintf("%s:%s", cl, ch)

	err := r.redisClient.Set(ctx, key, "", ttl).Err()
	if err != nil {
		log.Printf("error managing ttl: %v\n", err)
	}
	return err
}

func (r *heartbeatRepository) ListenMessages() {
	err := r.redisPubSub.Subscribe(context.Background(), "Heartbeat")
	if err != nil {
		log.Fatalf("subscribe error: %v", err)
	}

	ch := r.redisPubSub.Channel()
	for msg := range ch {

		if msg.Channel != "Heartbeat" {
			continue
		}

		fmt.Printf("Received message from %s: %s\n", msg.Channel, msg.Payload)

		var mesg Models.Message

		err := json.Unmarshal([]byte(msg.Payload), &mesg)
		if err != nil {
			fmt.Println("Error unmarshalling message: ", err)
			continue
		}

		if err := r.UserTTlManager(context.Background(), mesg.UserID, mesg.ChannelID); err != nil {
			log.Printf("error managing TTL for client %s: %v\n", msg.Payload, err)
		}
	}
}
