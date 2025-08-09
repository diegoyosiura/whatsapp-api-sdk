package domain

// MessageStatus represents delivery, read, failed and other lifecycle updates.
type MessageStatus struct {
	ID           string           `json:"id"`
	Status       string           `json:"status"`
	Timestamp    string           `json:"timestamp"`
	RecipientID  string           `json:"recipient_id"`
	Conversation *ConversationRef `json:"conversation,omitempty"`
	Pricing      *Pricing         `json:"pricing,omitempty"`
	Errors       []WebhookError   `json:"errors,omitempty"`
}
