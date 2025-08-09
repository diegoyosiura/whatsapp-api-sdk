package domain

// WebhookEvent is the root of the WhatsApp webhook payload.
// It mirrors the structure defined by the WhatsApp Cloud API Postman collection.
// Keep this package free of transport concerns and HTTP-specific details.
type WebhookEvent struct {
	Object string         `json:"object"`
	Entry  []WebhookEntry `json:"entry"`
}
