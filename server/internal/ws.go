package internal

import (
	"context"

	"github.com/SubhamMurarka/chat_app/models"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WsService interface {
	JoinRoomService(c context.Context, cl *models.Client)
	SendMessage(c context.Context, cl *models.Client, message models.Message)
	Broadcast(c context.Context, room string, msg models.Message)
	handleWebSocket(conn *websocket.Conn, c *gin.Context)
}
