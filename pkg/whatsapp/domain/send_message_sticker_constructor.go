package domain

import "github.com/diegoyosiura/whatsapp-sdk-go/pkg/errorsx"

type StickerMessage struct {
	Sticker *StickerBody `json:"sticker"`
}
type StickerBody struct {
	Id   *string `json:"id"`
	Link *string `json:"link"`
}

func NewSendStickerRequest(to, imageID, imageURL string) *SendMessage {
	rt := "individual"
	sb := &StickerBody{}

	if imageID != "" {
		sb.Id = &imageID
	} else {
		sb.Link = &imageURL
	}
	return &SendMessage{
		MessagingProduct: "whatsapp",
		RecipientType:    &rt,
		To:               to,
		Type:             "sticker",
		StickerMessage:   &StickerMessage{Sticker: sb},
	}
}

func NewSendContextStickerRequest(to, imageID, imageURL, targetMessage string) *SendMessage {
	rt := "individual"
	sb := &StickerBody{}

	if imageID != "" {
		sb.Id = &imageID
	} else {
		sb.Link = &imageURL
	}
	return &SendMessage{
		MessagingProduct: "whatsapp",
		RecipientType:    &rt,
		To:               to,
		Type:             "sticker",
		ContextMessage:   &ContextMessage{Context: &Context{MessageId: targetMessage}},
		StickerMessage:   &StickerMessage{Sticker: sb},
	}
}
func (s *SendMessage) validateStickerMessage() error {
	if s.Type != "sticker" {
		return &errorsx.ValidationError{Field: "Type", Reason: "type must be sticker", Op: "validateStickerMessage"}
	}

	if s.StickerMessage == nil {
		return &errorsx.ValidationError{Field: "StickerMessage", Reason: "sticker is nil", Op: "validateStickerMessage"}
	}

	if s.StickerMessage.Sticker == nil {
		return &errorsx.ValidationError{Field: "Sticker", Reason: "sticker is nil", Op: "validateStickerMessage"}
	}

	if *s.StickerMessage.Sticker.Id == "" && *s.StickerMessage.Sticker.Link == "" {
		return &errorsx.ValidationError{Field: "Sticker", Reason: "the reference must be a link or a Id, nothing received", Op: "validateStickerMessage"}
	}

	if *s.StickerMessage.Sticker.Id != "" && *s.StickerMessage.Sticker.Link != "" {
		return &errorsx.ValidationError{Field: "Sticker", Reason: "the reference must be a link or a Id, both received", Op: "validateStickerMessage"}
	}
	return nil
}
