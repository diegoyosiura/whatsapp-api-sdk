package services

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/diegoyosiura/whatsapp-sdk-go/pkg/whatsapp/domain"
)

func (s *MessagesService) SendImage(ctx context.Context, to, imageId, imageURL string) (*domain.MessageSendResponse, error) {
	payload := domain.NewSendImageRequest(to, imageId, imageURL)
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
func (s *MessagesService) SendImageReply(ctx context.Context, to, imageId, imageURL, targetMessageId string) (*domain.MessageSendResponse, error) {
	payload := domain.NewSendContextImageRequest(to, imageId, imageURL, targetMessageId)
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
