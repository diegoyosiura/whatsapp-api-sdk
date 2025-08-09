package domain

import (
	"encoding/json"
)

// ParseWebhookEvent decodes a raw JSON payload into WebhookEvent.
func ParseWebhookEvent(data []byte) (WebhookEvent, error) {
	var e WebhookEvent
	if err := json.Unmarshal(data, &e); err != nil {
		return WebhookEvent{}, err
	}
	return e, nil
}
