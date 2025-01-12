package main

import (
	"log"

	"github.com/SubhamMurarka/chat_app/Heartbeat/Redis"
	"github.com/SubhamMurarka/chat_app/Heartbeat/Repository"
)

func main() {
	redisdb, err := Redis.NewRedisDatabase()
	if err != nil {
		log.Fatalf("error connecting redis heartbeat : %v", err)
		return
	}
	defer redisdb.Close()
	redisObj := Repository.NewHeartBeatRepository(redisdb.GetClient(), redisdb.GetPubSub())
	redisObj.ListenMessages()
}
