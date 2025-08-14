package services

import (
	"context"
	"fmt"
	"io"

	"github.com/diegoyosiura/whatsapp-sdk-go/pkg/errorsx"
	"github.com/diegoyosiura/whatsapp-sdk-go/pkg/whatsapp/ports"
	"github.com/diegoyosiura/whatsapp-sdk-go/pkg/whatsapp/transport/graph"
)

// NewMessagesService creates a new MessagesService bound to a minimal client interface.
func NewMessagesService(c clientCore) *MessagesService { return &MessagesService{c: c} }

func (s *MessagesService) doRequest(ctx context.Context, payload ports.SendMessage) ([]byte, error) {
	base := s.c.BaseURL()
	if base == "" {
		base = graph.DefaultBaseURL
	}
	buf, err := payload.Buffer()
	if err != nil {
		return nil, &errorsx.ValidationError{Op: "doRequest", Field: "Buffer", Reason: err.Error()}
	}

	req, err := graph.NewSendHTTPRequest(ctx, base, s.c.Version(), s.c.PhoneNumberID(), buf, s.c.TokenProvider())
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
