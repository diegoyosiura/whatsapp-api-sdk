package domain

type WebhookContact struct {
	Profile *struct {
		Name string `json:"name"`
	} `json:"profile,omitempty"`
	WaID string `json:"wa_id"`
}
