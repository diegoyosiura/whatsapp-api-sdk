package graph

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/diegoyosiura/whatsapp-sdk-go/pkg/whatsapp/ports"
)

// NewSendHTTPRequest builds an *http.Request for POST /{PN-ID}/messages.
// The caller is responsible for executing it via the configured HTTPDoer and
// decoding the response into domain models.
func NewSendHTTPRequest(ctx context.Context, base, version, phoneNumberID string, payload *bytes.Buffer, tokenProvider ports.TokenProvider) (*http.Request, error) {
	url := MessagesEndpoint(base, version, phoneNumberID)

	token, err := tokenProvider.Token(ctx)
	if err != nil {
		return nil, fmt.Errorf("token: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	// Enable retries by making the body rewindable when supported by the runtime.
	b := payload.Bytes()
	req.GetBody = func() (r io.ReadCloser, err error) { return io.NopCloser(bytes.NewReader(b)), nil }
	return req, nil
}
