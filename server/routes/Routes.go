package routes

import (
	"net/http"

	"github.com/SubhamMurarka/chat_app/server/handlers"
	middleware "github.com/SubhamMurarka/chat_app/server/middlewares"
	"github.com/gin-gonic/gin"
)

var r *gin.Engine

func InitRouter(wsHandler *handlers.WsHandler) {
	r = gin.Default()
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"Health": "Ready"})
	})
	r.Use(middleware.Authenticate())
	r.GET("/ws/user/joinRoom", wsHandler.WebSocketHandler)
}

func Start(addr string) error {
	return r.Run(addr)
}
