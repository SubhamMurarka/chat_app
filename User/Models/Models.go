package Models

import (
	"database/sql"
	"time"
)

type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateUserReq struct {
	Username string `json:"username" validate:"required,min=3,max=20"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type CreateUserRes struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type LoginUserReq struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginUserRes struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Token    string `json:"token"`
}

type Room struct {
	RoomID   int64  `json:"roomid"`
	RoomName string `json:"roomname" validate:"required,min=3,max=50"`
	Typ      string `json:"typ" validate:"required,oneof=DM GROUP"`
}

type UserOffline struct {
	Action    string `json:"action"`
	UserID    int64  `json:"user_id"`
	ChannelID int64  `json:"channel_id"`
}

type Message struct {
	ID          int64     `json:"id,omitempty"`
	Content     string    `json:"content,omitempty"`
	UserID      int64     `json:"user_id,omitempty"`
	ChannelID   int64     `json:"channel_id,omitempty"`
	MessageType string    `json:"message_type" validate:"required,oneof=TEXT MEDIA"`
	MediaID     string    `json:"media_id,omitempty"`
	Created_At  time.Time `json:"created_at"`
}

type MessageOP struct {
	ID          int64          `json:"id"`
	ChannelID   int64          `json:"channel_id"`
	UserID      int64          `json:"user_id"`
	Content     sql.NullString `json:"content"`
	MediaID     sql.NullString `json:"media_id"`
	MessageType string         `json:"message_type"`
	Created_At  time.Time      `json:"created_at"`
}

type Members struct {
	UserID    int64     `json:"user_id" validate:"required"`
	ChannelID int64     `json:"channel_id" validate:"required"`
	Username  string    `json:"username,omitempty" validate:"max=20"`
	Online_At time.Time `json:"online_at"`
}
