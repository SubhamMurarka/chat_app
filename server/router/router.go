package router

import (
	auth "github.com/SubhamMurarka/chat_app/internal/auth_middleware"
	"github.com/SubhamMurarka/chat_app/internal/user"
	"github.com/SubhamMurarka/chat_app/internal/ws"
	"github.com/gin-gonic/gin"
)

var r *gin.Engine

func InitRouter(userHandler *user.Handler, wsHandler *ws.Handler) {
	r = gin.Default()
	// r.Static("/home", "/home/murarka/chat_app/frontend")

	r.Use(gin.Logger())
	r.POST("/signup", userHandler.CreateUser)
	r.POST("/login", userHandler.Login)

	r.Use(auth.Authenticate())
	r.POST("/ws/createRoom", wsHandler.CreateRoom)
	r.GET("/ws/joinRoom/:roomId", wsHandler.JoinRoom)
	r.GET("/ws/getRooms", wsHandler.GetRooms)
	r.GET("/ws/getClients/:roomId", wsHandler.GetClients)
	r.POST("/upload", wsHandler.ImageUpload)
}

func Start(addr string) error {
	return r.Run(addr)
}
