package main

import (
	"fmt"
	"log"

	"github.com/SubhamMurarka/chat_app/server/AbuseMasking"
	"github.com/SubhamMurarka/chat_app/server/config"
	"github.com/SubhamMurarka/chat_app/server/db"
	"github.com/SubhamMurarka/chat_app/server/handlers"
	"github.com/SubhamMurarka/chat_app/server/helpers"
	"github.com/SubhamMurarka/chat_app/server/repositories"
	"github.com/SubhamMurarka/chat_app/server/routes"
	"github.com/SubhamMurarka/chat_app/server/services"
)

// var serverID

func main() {
	// Initialize Redis
	redisClient, err := db.NewRedisDatabase()
	if err != nil {
		log.Fatalf("could not initialize Redis connection: %s", err)
	}

	//initialize conn package
	loc := helpers.NewLocation()

	//Initialize AbuseMasker
	words := AbuseMasking.Loadfile()
	fmt.Println(words)
	AbuseMasking.MakeTrie(words)

	// Initialize Repositories
	pubsubRepo := repositories.NewPubSubRepository(redisClient.GetPubSub(), redisClient.GetClient(), loc)

	wsSvc := services.NewwsService(pubsubRepo, loc)
	wsHandler := handlers.NewWsHandler(wsSvc)

	routes.InitRouter(wsHandler)

	addr := fmt.Sprintf("%s:%s", config.Config.ServerHost, config.Config.ServerPort)
	routes.Start(addr)
}
