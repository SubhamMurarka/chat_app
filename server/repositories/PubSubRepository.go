package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/SubhamMurarka/chat_app/models"
	"github.com/redis/go-redis/v9"
)

type PubSubRepository interface {
	SubscribeRoom(ctx context.Context, room string, callback models.MessageHandlerCallbackType)
	PublishMessage(ctx context.Context, room string, msg *models.Message)
}

type pubsubRepository struct {
	pubsubConnection       *redis.PubSub
	redisClient            *redis.Client
	roomsMutex             *sync.Mutex
	roomsSubscribed        map[string]bool
	once                   *sync.Once
	messageHandlerCallback models.MessageHandlerCallbackType
}

func NewPubSubRepository(PubSubConnection *redis.PubSub, RedisClient *redis.Client) PubSubRepository {
	return &pubsubRepository{
		pubsubConnection: PubSubConnection,
		redisClient:      RedisClient,
		roomsMutex:       &sync.Mutex{},
		roomsSubscribed:  make(map[string]bool),
		once:             &sync.Once{},
	}
}

func (r *pubsubRepository) SubscribeRoom(c context.Context, room string, callback models.MessageHandlerCallbackType) {
	r.roomsMutex.Lock()
	defer r.roomsMutex.Unlock()

	if !r.roomsSubscribed[room] {
		r.pubsubConnection.Subscribe(c, room)
		r.roomsSubscribed[room] = true

		r.messageHandlerCallback = callback

		r.once.Do(func() {
			go r.listenMessages()
		})
	}
}

func (r *pubsubRepository) listenMessages() {
	channel := r.pubsubConnection.Channel()

	for message := range channel {
		var msg models.Message
		err := json.Unmarshal([]byte(message.Payload), &msg)
		if err != nil {
			fmt.Printf("error unmarshalling messages")
			continue
		}

		if r.messageHandlerCallback != nil {
			r.messageHandlerCallback(message.Channel, &msg)
		}
	}
}

func (r *pubsubRepository) PublishMessage(c context.Context, room string, msg *models.Message) {

	message, err := json.Marshal(msg)
	if err != nil {
		fmt.Printf("error marshalling message")
		return
	}

	r.redisClient.Publish(c, room, message)
}
