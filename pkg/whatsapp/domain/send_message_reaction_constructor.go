package domain

import (
	"github.com/diegoyosiura/whatsapp-sdk-go/pkg/errorsx"
)

type ReactionMessage struct {
	Reaction *ReactionBody `json:"reaction"`
}

type ReactionBody struct {
	MessageID string `json:"message_id"`
	Emoji     string `json:"emoji"`
}

func NewSendReplyReaction(to, messageID, emoji string) *SendMessage {
	rt := "individual"
	return &SendMessage{
		MessagingProduct: "whatsapp",
		RecipientType:    &rt,
		To:               to,
		Type:             "reaction",
		ReactionMessage: &ReactionMessage{
			Reaction: &ReactionBody{
				MessageID: messageID,
				Emoji:     emoji,
			},
		},
	}
}

func (s *SendMessage) validateReplyMessage() *errorsx.ValidationError {
	if s.Type == "reaction" {
		return &errorsx.ValidationError{Op: "validateReplyMessage", Field: "Type", Reason: "must be reaction"}
	}
	if s.Reaction.Emoji == "" {
		return &errorsx.ValidationError{Op: "validateReplyMessage", Field: "Emoji", Reason: "empty"}
	}

	if s.Reaction.MessageID == "" {
		return &errorsx.ValidationError{Op: "validateReplyMessage", Field: "MessageID", Reason: "empty"}
	}
	return nil
}
