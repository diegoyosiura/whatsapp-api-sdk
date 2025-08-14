package services

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/diegoyosiura/whatsapp-sdk-go/pkg/whatsapp/domain"
)

// SendText sends a simple text message to a phone number in E.164 format.
// It validates inputs, builds the HTTP request via transport, executes the
// request using the Client, and decodes the response into domain types.
func (s *MessagesService) SendText(ctx context.Context, to, body string) (*domain.MessageSendResponse, error) {
	payload := domain.NewSendTextMessage(to, body)
	if err := payload.Validate(); err != nil {
		return nil, err
	}

	b, err := s.doRequest(ctx, payload)

	if err != nil {
		return nil, err
	}

	var out domain.MessageSendResponse
	if err = json.Unmarshal(b, &out); err != nil {
		return nil, fmt.Errorf("decode success response: %w", err)
	}
	return &out, nil
}

func (s *MessagesService) SendTextReply(ctx context.Context, to, body, targetMessageId string) (*domain.MessageSendResponse, error) {
	payload := domain.NewSendContextTextRequest(to, body, targetMessageId)
	if err := payload.Validate(); err != nil {
		return nil, err
	}

	b, err := s.doRequest(ctx, payload)

	if err != nil {
		return nil, err
	}

	var out domain.MessageSendResponse
	if err = json.Unmarshal(b, &out); err != nil {
		return nil, fmt.Errorf("decode success response: %w", err)
	}
	return &out, nil
}
