package main

import (
	"log"

	"github.com/SubhamMurarka/chat_app/db"
	"github.com/SubhamMurarka/chat_app/internal"
	"github.com/SubhamMurarka/chat_app/reddis"
	"github.com/SubhamMurarka/chat_app/router"
)

func main() {
	reddis.InitRedis()

	dbConn, err := db.NewDatabase()
	if err != nil {
		log.Fatalf("could not initialiaze database connection: %s", err)
	}

	Rep := internal.NewRepository(dbConn.GetDB())
	userSvc := internal.NewService(Rep)
	userHandler := internal.NewHandler(userSvc)

	wsSvc := internal.NewwsService(Rep)
	wsHandler := internal.NewWsHandler(wsSvc)

	router.InitRouter(userHandler, wsHandler)

	router.Start("0.0.0.0:8080")
}
