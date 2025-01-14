package models

type Message struct {
	ID          uint64 `json:"id"`
	Content     string `json:"content,omitempty"`
	Server      string `json:"server"`
	UserID      int64  `json:"user_id"`
	ChannelID   int64  `json:"channel_id"`
	MessageType string `json:"message_type"` //text, images etc
	EventType   string `json:"event_type"`   //joinroom, chat, heartbeat
	MediaID     string `json:"media_id,omitempty"`
}
