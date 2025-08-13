package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/diegoyosiura/whatsapp-sdk-go/pkg/errorsx"
	"github.com/diegoyosiura/whatsapp-sdk-go/pkg/whatsapp/domain"
	"github.com/diegoyosiura/whatsapp-sdk-go/pkg/whatsapp/transport/graph"
)

// NewMessagesService creates a new MessagesService bound to a minimal client interface.
func NewMessagesService(c clientCore) *MessagesService { return &MessagesService{c: c} }

// SendText sends a simple text message to a phone number in E.164 format.
// It validates inputs, builds the HTTP request via transport, executes the
// request using the Client, and decodes the response into domain types.
func (s *MessagesService) SendText(ctx context.Context, to, body string) (*domain.MessageSendResponse, error) {
	if to == "" {
		return nil, &errorsx.ValidationError{Op: "SendText", Field: "to", Reason: "empty"}
	}
	if body == "" {
		return nil, &errorsx.ValidationError{Op: "SendText", Field: "body", Reason: "empty"}
	}

	if !isE164(to) {
		return nil, &errorsx.ValidationError{Op: "SendText", Field: "to", Reason: "must be E.164 like +5511999999999"}
	}

	payload := domain.SendTextRequest{MessagingProduct: "whatsapp", To: to, Type: "text"}
	payload.Text.Body = body
	buf, err := payload.Buffer()
	if err != nil {
		return nil, &errorsx.ValidationError{Op: "SendText", Field: "body", Reason: err.Error()}
	}
	b, err := s.doRequest(ctx, buf)

	if err != nil {
		return nil, err
	}

	var out domain.MessageSendResponse
	if err = json.Unmarshal(b, &out); err != nil {
		return nil, fmt.Errorf("decode success response: %w", err)
	}
	return &out, nil
}

func (s *MessagesService) doRequest(ctx context.Context, payload *bytes.Buffer) ([]byte, error) {
	base := s.c.BaseURL()
	if base == "" {
		base = graph.DefaultBaseURL
	}

	req, err := graph.NewSendHTTPRequest(ctx, base, s.c.Version(), s.c.PhoneNumberID(), payload, s.c.TokenProvider())
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	resp, err := s.c.Do(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		ge := errorsx.TryParseGraphError(resp, b)
		if ge != nil {
			return nil, ge
		}
		return nil, errorsx.NewHTTPErrorFromResponse(resp, b)
	}

	return b, nil
}
func isE164(s string) bool {
	if len(s) < 4 || len(s) > 17 {
		return false
	}
	if s[0] != '+' {
		return false
	}
	for i := 1; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return true
}
