package main

import (
	"log"

	"github.com/SubhamMurarka/chat_app/db"
	"github.com/SubhamMurarka/chat_app/internal"
	"github.com/SubhamMurarka/chat_app/repositories"
	"github.com/SubhamMurarka/chat_app/router"
	"github.com/SubhamMurarka/chat_app/services"
)

func main() {
	// Initialize Redis
	redisClient, err := db.NewRedisDatabase()
	if err != nil {
		log.Fatalf("could not initialize Redis connection: %s", err)
	}

	// Initialize postgres
	dbConn, err := db.NewSQLDatabase()
	if err != nil {
		log.Fatalf("could not initialiaze database connection: %s", err)
	}

	// Initialize Repositories
	userRepo := repositories.NewUserRepository(dbConn.GetDB())
	roomRepo := repositories.NewRoomRepository(redisClient.GetClient())
	pubsubRepo := repositories.NewPubSubRepository(redisClient.GetPubSub(), redisClient.GetClient())

	userSvc := services.NewUserService(userRepo)
	userHandler := .NewHandler(userSvc)

	wsSvc := services.NewwsService(roomRepo, pubsubRepo)
	wsHandler := internal.NewWsHandler(wsSvc)

	router.InitRouter(userHandler, wsHandler)

	router.Start("0.0.0.0:8080")
}
