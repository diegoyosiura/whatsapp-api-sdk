package domain

type AudioMessage struct {
	Audio *AudioBody `json:"audio"`
}
type AudioBody struct {
	Id   *string `json:"id"`
	Link *string `json:"link"`
}

func NewSendAudioRequest(to, imageID, imageURL string) *SendMessage {
	rt := "individual"
	ab := &AudioBody{}

	if imageID != "" {
		ab.Id = &imageID
	} else {
		ab.Link = &imageURL
	}
	return &SendMessage{
		MessagingProduct: "whatsapp",
		RecipientType:    &rt,
		To:               to,
		Type:             "audio",
		AudioMessage:     &AudioMessage{Audio: ab},
	}
}

func NewSendContextAudioRequest(to, imageID, imageURL, targetMessage string) *SendMessage {
	rt := "individual"
	ab := &AudioBody{}

	if imageID != "" {
		ab.Id = &imageID
	} else {
		ab.Link = &imageURL
	}
	return &SendMessage{
		MessagingProduct: "whatsapp",
		RecipientType:    &rt,
		To:               to,
		Type:             "audio",
		ContextMessage:   &ContextMessage{Context: &Context{MessageId: targetMessage}},
		AudioMessage:     &AudioMessage{Audio: ab},
	}
}
func (s *SendMessage) validateAudioMessage() error {
	return nil
}
