package router

import (
	"github.com/SubhamMurarka/chat_app/internal"
	auth "github.com/SubhamMurarka/chat_app/internal/auth_middleware"
	"github.com/gin-gonic/gin"
)

var r *gin.Engine

func InitRouter(userHandler *internal.Handler, wsHandler *internal.WsHandler) {
	r = gin.Default()
	// r.Static("/home", "/home/murarka/chat_app/frontend")

	r.POST("/signup", userHandler.CreateUser)
	r.POST("/login", userHandler.Login)

	r.Use(auth.Authenticate())
	r.POST("/upload", wsHandler.ImageUpload)
	r.GET("/ws/joinRoom", wsHandler.WebSocketHandler)
	// r.GET("/ws/getRooms", wsHandler.GetRooms)
	// r.GET("/ws/getClients/:roomId", wsHandler.GetClients)
}

func Start(addr string) error {
	return r.Run(addr)
}
