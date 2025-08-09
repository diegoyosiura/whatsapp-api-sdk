package domain

// WebhookValue captures the union of message and status updates that the API may send.
// Only the most common, stable fields are represented here to keep the domain surface stable.
// Less common fields can be added incrementally without breaking changes.
type WebhookValue struct {
	MessagingProduct string           `json:"messaging_product"`
	Metadata         *WebhookMetadata `json:"metadata,omitempty"`
	Contacts         []WebhookContact `json:"contacts,omitempty"`
	Messages         []InboundMessage `json:"messages,omitempty"`
	Statuses         []MessageStatus  `json:"statuses,omitempty"`
	Errors           []WebhookError   `json:"errors,omitempty"`
}
