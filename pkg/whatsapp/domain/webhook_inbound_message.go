package domain

// InboundMessage represents a message received by the business.
// Message types are additive; only the fields for the concrete message type will be present.
// Extend with other message types as needed using optional pointers.
type InboundMessage struct {
	ID          string             `json:"id"`
	From        string             `json:"from"`
	Timestamp   string             `json:"timestamp"`
	Type        string             `json:"type"`
	Text        *MessageText       `json:"text,omitempty"`
	Image       *MediaObject       `json:"image,omitempty"`
	Document    *MediaObject       `json:"document,omitempty"`
	Audio       *MediaObject       `json:"audio,omitempty"`
	Video       *MediaObject       `json:"video,omitempty"`
	Sticker     *MediaObject       `json:"sticker,omitempty"`
	Interactive *InteractiveObject `json:"interactive,omitempty"`
	Context     *MessageContext    `json:"context,omitempty"`
}
