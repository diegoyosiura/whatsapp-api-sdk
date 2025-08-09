package domain

type WebhookChange struct {
	Field string       `json:"field"`
	Value WebhookValue `json:"value"`
}
