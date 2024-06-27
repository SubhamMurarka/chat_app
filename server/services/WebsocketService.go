package services

import (
	"context"
	"fmt"
	"log"

	"github.com/SubhamMurarka/chat_app/helpers"
	"github.com/SubhamMurarka/chat_app/models"
	"github.com/SubhamMurarka/chat_app/repositories"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WsService interface {
	JoinRoomService(c context.Context, cl *models.Client)
	SendMessage(c context.Context, cl *models.Client, message models.Message)
	Broadcast(c context.Context, room string, msg models.Message)
	handleWebSocket(conn *websocket.Conn, c *gin.Context)
}

type wsService struct {
	roomrepo   repositories.RoomRepository
	pubsubrepo repositories.PubSubRepository
	// timeout time.Duration
}

func NewwsService(roomRepository repositories.RoomRepository, pubsubRepository repositories.PubSubRepository) WsService {
	return &wsService{
		roomrepo:   roomRepository,
		pubsubrepo: pubsubRepository,
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

	helpers.AddConnection(clientID, conn)

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
	err := s.roomrepo.AddUserToRoomRedis(c, room, cl)
	if err != nil {
		cl.Conn.WriteJSON(err.Error())
		return
	}

	s.pubsubrepo.SubscribeRoom(c, room, func(room string, message *models.Message) {
		s.Broadcast(c, room, *message)
	})

	message := models.Message{
		RoomID:      cl.RoomID,
		Username:    cl.Username,
		UserID:      cl.ID,
		MessageType: "join_room",
	}

	s.pubsubrepo.PublishMessage(c, room, &message)

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
	users := s.roomrepo.GetAllMembersRedis(c, room)

	for _, userID := range users {

		if conn, ok := helpers.GetConnection(userID); ok {
			if err := conn.WriteJSON(msg); err != nil {
				// handle error
				fmt.Printf("error in writing to connection")
				conn.Close()
				helpers.RemoveConnection(userID)
			}
		}
	}
}

func (s *wsService) SendMessage(c context.Context, cl *models.Client, message models.Message) {
	// ctx, cancel := context.WithTimeout(c, s.timeout)
	// defer cancel()

	s.pubsubrepo.PublishMessage(c, cl.RoomID, &message)
}
