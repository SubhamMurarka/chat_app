package reddis

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/SubhamMurarka/chat_app/models"
)

var roomsMutex = &sync.Mutex{}
var roomsSubscribed = make(map[string]bool)

var once = &sync.Once{}

var messageHandlerCallback models.MessageHandlerCallbackType

func SubscribeRoom(c context.Context, room string, callback models.MessageHandlerCallbackType) {
	roomsMutex.Lock()
	defer roomsMutex.Unlock()

	if !roomsSubscribed[room] {
		// ctx, cancel := context.WithTimeout(c, s)
		// defer cancel()

		PubSubConnection.Subscribe(c, room)
		roomsSubscribed[room] = true

		messageHandlerCallback = callback

		// firing listenformessages once for every server
		once.Do(func() {
			go listenMessages()
		})
	}
}

func listenMessages() {
	// getting all the channels troom stringo check for messages
	channel := PubSubConnection.Channel()

	// checking for messages for each channel/room
	// conn.WriteJSON("hey go is fired")
	for message := range channel {
		var msg models.Message
		err := json.Unmarshal([]byte(message.Payload), &msg)

		if err != nil {
			// Todo error handling
			fmt.Printf("error unmarshalling messages")
			continue
		}

		if messageHandlerCallback != nil {
			messageHandlerCallback(message.Channel, &msg)
		}

	}
}

func PublishMessage(c context.Context, room string, msg *models.Message) {
	// ctx, cancel := context.WithTimeout(c, s)
	// defer cancel()

	message, err := json.Marshal(msg)
	if err != nil {
		// Todo error handling
		fmt.Printf("error marshalling message")
		return
	}

	RedisClient.Publish(c, room, message)
}
