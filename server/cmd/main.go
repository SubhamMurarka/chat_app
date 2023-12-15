package main

import (
	"log"

	"github.com/SubhamMurarka/chat_app/db"
)

func main() {
	_, err := db.NewDatabase()
	if err != nil {
		log.Fatalf("could not initialiaze database connection: %s", err)
	}
}
