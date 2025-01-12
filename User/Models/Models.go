package Models

import "time"

type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email" `
	Password string `json:"password" `
}

type CreateUserReq struct {
	Username string `json:"username" `
	Email    string `json:"email" `
	Password string `json:"password" `
}

type CreateUserRes struct {
	ID       string `json:"id" `
	Username string `json:"username" `
	Email    string `json:"email" `
}

type LoginUserReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginUserRes struct {
	ID       string `json:"id" `
	Username string `json:"username"`
	Token    string `json:"token"`
}

type Room struct {
	RoomID   int64  `json:"roomid" `
	RoomName string `json:"roomname" `
	Typ      string `json:"typ"`
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
	MessageType string    `json:"message_type" validate:"required"` //text, images etc
	MediaID     string    `json:"media_id,omitempty"`
	Created_At  time.Time `json:"created_at"`
}

type Members struct {
	UserID    int64     `json:"user_id"`
	ChannelID int64     `json:"channel_id"`
	Username  string    `json:"username,omitempty"`
	Online_At time.Time `json:"online_at"`
	Page      int       `json:"page"`
}
