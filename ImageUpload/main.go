package main

import (
	"fmt"
	"log"

	"github.com/SubhamMurarka/chat_app/Image/Config"
	"github.com/SubhamMurarka/chat_app/Image/S3"
	"github.com/SubhamMurarka/chat_app/Image/redis"
	util "github.com/SubhamMurarka/chat_app/Image/utils"
	"github.com/gin-gonic/gin"
)

func main() {

	err := redis.NewRedisDatabase()
	if err != nil {
		log.Fatalf("Error connecting redis")
		return
	}

	r := gin.Default()

	r.Use(util.Authenticate())
	r.GET("/img/url", S3.ConnectS3)

	addr := fmt.Sprintf("%s:%s", Config.Conf.AppHost, Config.Conf.AppPort)

	r.Run(addr)
}
