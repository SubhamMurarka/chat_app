package main

import (
	"log"

	"github.com/SubhamMurarka/chat_app/db"
	"github.com/SubhamMurarka/chat_app/internal/user"
	"github.com/SubhamMurarka/chat_app/internal/ws"
	"github.com/SubhamMurarka/chat_app/router"
)

func main() {
	dbConn, err := db.NewDatabase()
	if err != nil {
		log.Fatalf("could not initialiaze database connection: %s", err)
	}

	userRep := user.NewRepository(dbConn.GetDB())
	userSvc := user.NewService(userRep)
	userHandler := user.NewHandler(userSvc)

	hub := ws.NewHub()
	wsHandler := ws.NewHandler(hub)
	go hub.Run()

	router.InitRouter(userHandler, wsHandler)

	router.Start("0.0.0.0:8080")
}
