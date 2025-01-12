package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/SubhamMurarka/chat_app/server/config"
	"github.com/SubhamMurarka/chat_app/server/helpers"
	"github.com/SubhamMurarka/chat_app/server/models"
	"github.com/redis/go-redis/v9"
)

type PubSubRepository interface {
	SubscribeRoom(ctx context.Context, room string, channel_id int64, callback models.MessageHandlerCallbackType)
	PublishMessage(ctx context.Context, room string, channel_id int64, msg *models.Message)
}

type pubsubRepository struct {
	pubsubConnection       *redis.PubSub
	redisClient            *redis.Client
	roomsMutex             *sync.Mutex
	roomsSubscribed        map[int64]bool
	messageHandlerCallback models.MessageHandlerCallbackType
	once                   *sync.Once
	location               *helpers.Location
}

func NewPubSubRepository(PubSubConnection *redis.PubSub, RedisClient *redis.Client, loc *helpers.Location) PubSubRepository {
	return &pubsubRepository{
		pubsubConnection: PubSubConnection,
		redisClient:      RedisClient,
		roomsMutex:       &sync.Mutex{},
		roomsSubscribed:  make(map[int64]bool),
		once:             &sync.Once{},
		location:         loc,
	}
}

func (r *pubsubRepository) SubscribeRoom(c context.Context, room string, channel_id int64, callback models.MessageHandlerCallbackType) {
	r.roomsMutex.Lock()
	defer r.roomsMutex.Unlock()

	if !r.roomsSubscribed[channel_id] {
		r.pubsubConnection.Subscribe(c, room)
		r.roomsSubscribed[channel_id] = true
		r.messageHandlerCallback = callback

		r.once.Do(func() {
			go r.listenMessages()
		})
	}
}

func (r *pubsubRepository) listenMessages() {

	channel := r.pubsubConnection.Channel()

	for msg := range channel {

		if msg.Channel == "__keyevent@0__:expired" {
			fmt.Println(msg.Payload)
			keys := strings.Split(msg.Payload, ":")
			fmt.Println(keys, "checking key parsing")
			if len(keys) == 2 {
				clientID, _ := strconv.ParseInt(keys[0], 10, 64)
				channelID, _ := strconv.ParseInt(keys[1], 10, 64)
				_, exists := r.location.FetchUserConn(channelID, clientID)
				fmt.Println(clientID, channelID, exists, "checking exists")
				if exists {
					fmt.Println(clientID, channelID, "checking exists")
					r.location.RemoveUserFromRoom(channelID, clientID)
				}
			}
		} else {
			var mesg models.Message
			err := json.Unmarshal([]byte(msg.Payload), &mesg)
			if err != nil {
				fmt.Printf("error unmarshalling messages")
				continue
			}

			if r.messageHandlerCallback != nil && mesg.Server != config.Config.ServerID {
				fmt.Println("broadcasting again", mesg.Server, config.Config.ServerID)
				r.messageHandlerCallback(msg.Channel, mesg.ChannelID, &mesg)
			}
		}
	}

}

func (r *pubsubRepository) PublishMessage(c context.Context, room string, channel_id int64, msg *models.Message) {
	message, err := json.Marshal(msg)
	if err != nil {
		fmt.Println("Error marshalling message: ", err)
		return
	}

	r.redisClient.Publish(c, room, message)
}
