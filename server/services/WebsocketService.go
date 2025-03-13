package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/SubhamMurarka/chat_app/server/AbuseMasking"
	"github.com/SubhamMurarka/chat_app/server/config"
	"github.com/SubhamMurarka/chat_app/server/helpers"
	"github.com/SubhamMurarka/chat_app/server/kafka"
	"github.com/SubhamMurarka/chat_app/server/models"
	"github.com/SubhamMurarka/chat_app/server/repositories"
	"github.com/SubhamMurarka/chat_app/server/util"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/websocket"
)

var LLMID uint64

type WsService interface {
	JoinRoomService(c context.Context, cl *models.Client, message *models.Message)
	SendMessage(c context.Context, cl *models.Client, message *models.Message)
	Broadcast(c context.Context, room string, channelID int64, msg *models.Message)
	HandleWebSocket(c *gin.Context, cl *models.Client)
	HeartBeat(c context.Context, cl *models.Client, msg *models.Message)
}

type wsService struct {
	pubsubrepo repositories.PubSubRepository
	location   *helpers.Location
}

func NewwsService(pubsubRepository repositories.PubSubRepository, loc *helpers.Location) WsService {
	return &wsService{
		pubsubrepo: pubsubRepository,
		location:   loc,
	}
}

func (s *wsService) HandleWebSocket(c *gin.Context, cl *models.Client) {
	defer cl.Conn.Close() // Ensure the connection is closed when done

	server := config.Config.ServerID
	log.Printf("Server ID: %s", server)

	joinMessage := &models.Message{
		Server:      server,
		UserID:      cl.ClientID,
		ChannelID:   cl.ChannelID,
		EventType:   "join_room",
		Content:     fmt.Sprintf("%v joined", cl.UserName),
		MessageType: "TEXT",
	}
	s.JoinRoomService(c, cl, joinMessage)

	for {
		var message models.Message
		message.Server = server
		message.ChannelID = cl.ChannelID
		message.UserID = cl.ClientID

		var validate = validator.New()

		log.Println("Waiting for client message...")
		if err := cl.Conn.ReadJSON(&message); err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}

		if validationerr := validate.Struct(&message); validationerr != nil {
			cl.Conn.WriteJSON(validationerr)
			continue
		}

		log.Printf("Received message: %v", message.Content)
		s.handleMessage(c, cl, &message)
	}
}

func (s *wsService) handleMessage(c context.Context, cl *models.Client, message *models.Message) {
	switch message.EventType {
	case "chat":
		message.Content = AbuseMasking.Filter(message.Content)
		s.SendMessage(c, cl, message)
	case "heartbeat":
		s.HeartBeat(c, cl, message)
	default:
		log.Printf("Unknown event type: %v", message.EventType)
	}
}

func (s *wsService) JoinRoomService(c context.Context, cl *models.Client, message *models.Message) {
	s.location.AddUserToRoom(cl.ChannelID, cl.ClientID, cl.Conn)

	// Subscribe to room
	s.pubsubrepo.SubscribeRoom(c, cl.RoomName, cl.ChannelID, func(room string, channelID int64, msg *models.Message) {
		s.Broadcast(c, room, channelID, msg)
	})

	// Broadcast join message
	s.Broadcast(c, cl.RoomName, cl.ChannelID, message)
	s.pubsubrepo.PublishMessage(c, cl.RoomName, cl.ChannelID, message)

	response := map[string]interface{}{
		"type":    "join_room",
		"room":    cl.RoomName,
		"success": true,
	}
	if err := cl.Conn.WriteJSON(response); err != nil {
		log.Printf("Error sending join_room response: %v", err)
	}

	s.HeartBeat(c, cl, message) // Send initial heartbeat
}

func (s *wsService) Broadcast(c context.Context, room string, channelID int64, msg *models.Message) {
	users, exists := s.location.FetchUsersInRoom(channelID)
	if !exists || len(users) == 0 {
		log.Printf("No users found in room: %s", room)
		return
	}

	messageBytes, err := encodeMessage(msg)
	if err != nil {
		log.Printf("Error encoding message: %v", err)
		return
	}

	for userID, conn := range users {
		if err := conn.WriteMessage(websocket.TextMessage, messageBytes); err != nil {
			log.Printf("Error writing to user %d: %v", userID, err)
			s.location.RemoveUserFromRoom(channelID, userID)
			conn.Close()
		}
	}
}

func (s *wsService) SendMessage(c context.Context, cl *models.Client, message *models.Message) {
	message.ID, _ = helpers.GenerateID()
	fmt.Println("checking message", message)
	s.Broadcast(c, cl.RoomName, cl.ChannelID, message)
	s.pubsubrepo.PublishMessage(c, cl.RoomName, cl.ChannelID, message)
	kafka.ProduceToKafka(*message)
	if strings.HasPrefix(message.Content, "@superchat") {
		response := util.GetResponseFromLLM(message.Content)
		ID, _ := helpers.GenerateID()
		log.Printf("Response from LLM")
		llmMessage := &models.Message{
			ID:        ID,
			Content:   response,
			Server:    config.Config.ServerID,
			ChannelID: message.ChannelID,
			UserID:    123,
			EventType: "chat",
		}
		s.Broadcast(c, cl.RoomName, cl.ChannelID, llmMessage)
		s.pubsubrepo.PublishMessage(c, cl.RoomName, cl.ChannelID, llmMessage)
		kafka.ProduceToKafka(*llmMessage)
	}
}

func (s *wsService) HeartBeat(c context.Context, cl *models.Client, msg *models.Message) {
	msg.Content = "PONG"
	response := map[string]interface{}{
		"type":    "heartbeat",
		"message": "PONG",
	}
	if err := cl.Conn.WriteJSON(response); err != nil {
		log.Printf("Error sending heartbeat: %v", err)
	}
	// Optionally publish heartbeat messages to the pubsub
	s.pubsubrepo.PublishMessage(c, "Heartbeat", -123, msg)
}

func encodeMessage(msg *models.Message) ([]byte, error) {
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(msg); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
