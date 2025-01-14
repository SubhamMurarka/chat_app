package models

import (
	"github.com/gorilla/websocket"
)

type Client struct {
	Conn      *websocket.Conn
	ClientID  int64  `json:"id"`
	ChannelID int64  `json:"roomid"`
	UserName  string `json:"username"`
	RoomName  string `json:"roomname"`
}

type Message struct {
	ID          uint64 `json:"id,omitempty"`
	Content     string `json:"content,omitempty"`
	Server      string `json:"server"`
	UserID      int64  `json:"user_id,omitempty"`
	ChannelID   int64  `json:"channel_id,omitempty"`
	MessageType string `json:"message_type" validate:"required,oneof=TEXT MEDIA"`   //text, images etc
	EventType   string `json:"event_type" validate:"required,oneof=chat heartbeat"` //joinroom, chat, heartbeat
	MediaID     string `json:"media_id,omitempty"`
}

type Member struct {
	conn      *websocket.Conn
	channelid int64
}

type MessageHandlerCallbackType func(room string, channel_id int64, message *Message)
