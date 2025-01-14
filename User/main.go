package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/SubhamMurarka/chat_app/User/Config"
	db "github.com/SubhamMurarka/chat_app/User/DB"
	Handler "github.com/SubhamMurarka/chat_app/User/Handlers"
	middleware "github.com/SubhamMurarka/chat_app/User/Middleware"
	Repositories "github.com/SubhamMurarka/chat_app/User/Repository"
	"github.com/SubhamMurarka/chat_app/User/Services"
	"github.com/gin-gonic/gin"
)

func main() {
	conn, err := db.NewSQLDatabase()
	if err != nil {
		log.Fatal("Error getting Database %w", err)
	}

	connpubsub, err := db.NewRedisDatabase()
	if err != nil {
		log.Fatal("Error getting Database %w", err)
	}

	connmsg, err := db.NewMsgDatabase()
	if err != nil {
		log.Fatal("Error getting Database %w", err)
	}

	userRepo := Repositories.NewUserRepository(conn.GetDB())

	userSvc := Services.NewUserService(userRepo)

	roomRepo := Repositories.NewRoomRepository(conn.GetDB())

	msgRepo := Repositories.NewRepo(connmsg.GetDB())

	roomSvc := Services.NewRoomService(roomRepo)

	pubsubsvc := Services.NewPubSubService(connpubsub.GetPubSub(), connpubsub.GetClient(), roomRepo)

	msgSvc := Services.NewMsgService(msgRepo)

	Handler := Handler.NewUserHandler(userSvc, roomSvc, msgSvc)

	pubsubsvc.SubscribeRoom(context.Background())

	go pubsubsvc.ListenMessages(context.Background())

	r := gin.Default()

	//routes
	r.POST("/user/signup", Handler.CreateUser)
	r.POST("/user/login", Handler.Login)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"Health": "Ready"})
	})

	r.Use(middleware.Authenticate())
	r.GET("/user/joinroom", Handler.JoinRoom)
	r.GET("/user/getrooms", Handler.GetAllUserRoom)
	r.POST("/user/createroom", Handler.CreateRoom)
	r.GET("/user/members", Handler.GetAllMembers)
	r.GET("/user/history", Handler.GetMessages)

	addr := fmt.Sprintf("%s:%s", Config.Config.ServerHost, Config.Config.ServerPort)
	fmt.Println("port : ", Config.Config.ServerPort)
	r.Run(addr)
}
