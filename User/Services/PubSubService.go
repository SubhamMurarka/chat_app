package Services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/SubhamMurarka/chat_app/User/Models"
	"github.com/SubhamMurarka/chat_app/User/Repository"
	"github.com/SubhamMurarka/chat_app/User/Util"
	"github.com/redis/go-redis/v9"
)

type PubSubService interface {
	SubscribeRoom(c context.Context) error
	ListenMessages(c context.Context)
}

type pubsubService struct {
	pubsub     *redis.PubSub
	client     *redis.Client
	connection Repository.RoomRepository
}

func NewPubSubService(ps *redis.PubSub, cl *redis.Client, conn Repository.RoomRepository) PubSubService {
	return &pubsubService{
		pubsub:     ps,
		client:     cl,
		connection: conn,
	}
}

func (ps *pubsubService) SubscribeRoom(c context.Context) error {
	err := ps.pubsub.Subscribe(c, "__keyevent@0__:expired")
	if err != nil {
		log.Printf("Error subscribing expired keys channel : %w", err)
		return Util.ErrInternal
	}
	return nil
}

func (ps *pubsubService) ListenMessages(c context.Context) {
	ch := ps.pubsub.Channel()

	for val := range ch {
		if val.Channel != "__keyevent@0__:expired" {
			continue
		}

		var offline Models.UserOffline

		err := json.Unmarshal([]byte(val.Payload), &offline)
		if err != nil {
			fmt.Println("Error unmarshalling message: ", err)
			continue
		}

		err = ps.connection.UpdateStatus(c, offline.ChannelID, offline.UserID)
		if err != nil {
			log.Printf("Error updating online status for user %d of room %d", offline.UserID, offline.ChannelID)
		}
	}
}
