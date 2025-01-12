package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/SubhamMurarka/chat_app/server/models"
	"github.com/SubhamMurarka/chat_app/server/services"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WsHandler struct {
	WebsocketService services.WsService
}

func NewWsHandler(s services.WsService) *WsHandler {
	return &WsHandler{
		WebsocketService: s,
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow all origins for simplicity; update as per security needs
		return true
	},
}

func (h *WsHandler) WebSocketHandler(c *gin.Context) {
	// Room ID and User validation
	roomid := c.DefaultQuery("roomid", "")
	if roomid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Room parameter is required"})
		return
	}

	roomname := c.DefaultQuery("roomname", "")
	if roomname == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Room parameter is required"})
		return
	}

	clientID := c.GetString("userid")
	ClID, _ := strconv.ParseInt(clientID, 10, 64)

	clientName := c.GetString("username")

	channel_id, err := strconv.ParseInt(roomid, 10, 64)

	if err != nil {
		log.Printf("Error converting to string : %w", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	client := &models.Client{
		ClientID:  ClID,
		ChannelID: channel_id,
		UserName:  clientName,
		RoomName:  roomname,
	}

	// Initiate WebSocket upgrade
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Error upgrading to WebSocket:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upgrade to WebSocket"})
		return
	}

	client.Conn = conn

	// Call WebSocket service to handle the connection and the join-room logic
	h.WebsocketService.HandleWebSocket(c, client)
}
