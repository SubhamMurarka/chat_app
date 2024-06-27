package models

import "github.com/gorilla/websocket"

type User struct {
	ID       int64  `json:"id" db:"id"`
	Username string `json:"username" db:"username"`
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`
}

type CreateUserReq struct {
	Username string `json:"username" db:"username"`
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`
}

type CreateUserRes struct {
	ID       string `json:"id" db:"id"`
	Username string `json:"username" db:"username"`
	Email    string `json:"email" db:"email"`
}

type LoginUserReq struct {
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`
}

type LoginUserRes struct {
	ID       string `json:"id" db:"id"`
	Username string `json:"username" db:"username"`
	Token    string `json:"token"`
}

type Client struct {
	Conn     *websocket.Conn
	Message  *Message
	ID       string `json:"id"`
	RoomID   string `json:"roomid"`
	Username string `json:"username"`
}

type Message struct {
	Content     string   `json:"content"`
	RoomID      string   `json:"roomId"`
	Username    string   `json:"username"`
	UserID      string   `json:"userid"`
	Server      string   `json:"server"`
	MessageType string   `json:"messagetype"`
	ImageUrl    []string `json:"imageurl,omitempty"`
}

type MessageHandlerCallbackType func(room string, message *Message)
