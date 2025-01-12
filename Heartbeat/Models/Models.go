package Models

type Message struct {
	UserID    int64  `json:"user_id"`
	ChannelID int64  `json:"channel_id,omitempty"`
	EventType string `json:"even_type"`
}
