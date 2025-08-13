package domain

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// SendTextRequest is a minimal shape for sending a text message.
// The complete shape will live in the domain layer; here we only need a
// transport-level representation for the request body.
type SendTextRequest struct {
	MessagingProduct string `json:"messaging_product"`
	To               string `json:"to"`
	Type             string `json:"type"`
	Text             struct {
		Body string `json:"body"`
	} `json:"text"`
}

func (str *SendTextRequest) Buffer() (*bytes.Buffer, error) {
	buf := &bytes.Buffer{}
	if err := json.NewEncoder(buf).Encode(str); err != nil {
		return nil, fmt.Errorf("encode request: %w", err)
	}

	return buf, nil
}
