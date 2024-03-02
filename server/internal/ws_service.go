package internal

import (
	"context"
	"fmt"
	"log"

	"github.com/SubhamMurarka/chat_app/models"
	"github.com/SubhamMurarka/chat_app/reddis"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type wsService struct {
	Repository
	// timeout time.Duration
}

func NewwsService(repository Repository) WsService {
	return &wsService{
		repository,
		// time.Duration(8) * time.Hour,
	}
}

func (s *wsService) handleWebSocket(conn *websocket.Conn, c *gin.Context) {
	clientID := c.GetString("userid")
	clientName := c.GetString("username")

	cl := &models.Client{
		Conn:     conn,
		ID:       clientID,
		Username: clientName,
	}

	AddConnection(clientID, conn)

	for {
		var message models.Message

		err := conn.ReadJSON(&message)
		if err != nil {
			conn.WriteJSON(err.Error())
			return
		}

		message.UserID = clientID
		message.Username = clientName
		cl.RoomID = message.RoomID

		switch message.MessageType {
		case "join_room":
			s.JoinRoomService(context.Background(), cl)
		case "chat":
			s.SendMessage(context.Background(), cl, message)
		}
	}
}

func (s *wsService) JoinRoomService(c context.Context, cl *models.Client) {
	// ctx, cancel := context.WithTimeout(c, s.timeout)
	// defer cancel()

	room := cl.RoomID
	err := s.Repository.AddUserToRoomRedis(c, room, cl)
	if err != nil {
		cl.Conn.WriteJSON(err.Error())
		return
	}

	reddis.SubscribeRoom(c, room, func(room string, message *models.Message) {
		s.Broadcast(c, room, *message)
	})

	message := models.Message{
		RoomID:      cl.RoomID,
		Username:    cl.Username,
		UserID:      cl.ID,
		MessageType: "join_room",
	}

	reddis.PublishMessage(c, room, &message)

	response := map[string]interface{}{
		"type":    "join_room",
		"room":    room,
		"success": true,
	}

	if err := cl.Conn.WriteJSON(response); err != nil {
		log.Printf("%s", err)
	}
}

func (s *wsService) Broadcast(c context.Context, room string, msg models.Message) {
	// ctx, cancel := context.WithTimeout(c, s.timeout)
	// defer cancel()
	users := s.Repository.GetAllMembersRedis(c, room)

	for _, userID := range users {

		if conn, ok := GetConnection(userID); ok {
			if err := conn.WriteJSON(msg); err != nil {
				// handle error
				fmt.Printf("error in writing to connection")
				conn.Close()
				RemoveConnection(userID)
			}
		}
	}
}

func (s *wsService) SendMessage(c context.Context, cl *models.Client, message models.Message) {
	// ctx, cancel := context.WithTimeout(c, s.timeout)
	// defer cancel()

	reddis.PublishMessage(c, cl.RoomID, &message)
}
