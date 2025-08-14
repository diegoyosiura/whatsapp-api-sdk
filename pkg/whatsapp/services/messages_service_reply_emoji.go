package services

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/diegoyosiura/whatsapp-sdk-go/pkg/whatsapp/domain"
)

func (s *MessagesService) SendEmojiReply(ctx context.Context, to, targetMessageId, emoji string) (*domain.MessageSendResponse, error) {
	payload := domain.NewSendReplyReaction(to, targetMessageId, emoji)
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
