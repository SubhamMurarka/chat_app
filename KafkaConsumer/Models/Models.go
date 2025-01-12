package models

import "time"

type Message struct {
	ID          int64     `json:"id,omitempty"`
	Content     string    `json:"content,omitempty"`
	UserID      int64     `json:"user_id,omitempty"`
	ChannelID   int64     `json:"channel_id,omitempty"`
	MessageType string    `json:"message_type" validate:"required"` //text, images etc
	MediaID     string    `json:"media_id,omitempty"`
	Created_At  time.Time `json:"created_at"`
}
