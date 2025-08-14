package domain

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
	return nil
}
