package graph

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/diegoyosiura/whatsapp-sdk-go/pkg/whatsapp/ports"
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

// NewSendTextHTTPRequest builds an *http.Request for POST /{PN-ID}/messages.
// The caller is responsible for executing it via the configured HTTPDoer and
// decoding the response into domain models.
func NewSendTextHTTPRequest(ctx context.Context, base, version, phoneNumberID string, payload SendTextRequest, tokenProvider ports.TokenProvider) (*http.Request, error) {
	url := MessagesEndpoint(base, version, phoneNumberID)
	buf := &bytes.Buffer{}
	if err := json.NewEncoder(buf).Encode(payload); err != nil {
		return nil, fmt.Errorf("encode request: %w", err)
	}

	token, err := tokenProvider.Token(ctx)
	if err != nil {
		return nil, fmt.Errorf("token: %w", err)
	}
	
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	// Enable retries by making the body rewindable when supported by the runtime.
	b := buf.Bytes()
	req.GetBody = func() (r io.ReadCloser, err error) { return io.NopCloser(bytes.NewReader(b)), nil }
	return req, nil
}
