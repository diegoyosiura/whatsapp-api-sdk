package domain

type ConversationRef struct {
	ID     string  `json:"id"`
	Origin *Origin `json:"origin,omitempty"`
}
