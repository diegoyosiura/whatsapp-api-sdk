package domain

import (
	"github.com/diegoyosiura/whatsapp-sdk-go/pkg/errorsx"
	"github.com/diegoyosiura/whatsapp-sdk-go/pkg/utils"
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
	if s.Reaction.Emoji == "" {
		return &errorsx.ValidationError{Op: "SendReplyReaction", Field: "Emoji", Reason: "empty"}
	}

	if s.Reaction.MessageID == "" {
		return &errorsx.ValidationError{Op: "SendReplyReaction", Field: "MessageID", Reason: "empty"}
	}

	if !utils.IsE164(s.To) {
		return &errorsx.ValidationError{Op: "SendReplyReaction", Field: "to", Reason: "must be E.164 like +5511999999999"}
	}
	return nil
}
