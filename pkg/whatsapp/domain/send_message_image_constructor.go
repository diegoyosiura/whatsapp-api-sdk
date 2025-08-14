package domain

type ImageMessage struct {
	Image *ImageBody `json:"image"`
}
type ImageBody struct {
	Id   *string `json:"id"`
	Link *string `json:"link"`
}

func NewSendImageRequest(to, imageID, imageURL string) *SendMessage {
	rt := "individual"
	ib := &ImageBody{}

	if imageID != "" {
		ib.Id = &imageID
	} else {
		ib.Link = &imageURL
	}
	return &SendMessage{
		MessagingProduct: "whatsapp",
		RecipientType:    &rt,
		To:               to,
		Type:             "image",
		ImageMessage:     &ImageMessage{Image: ib},
	}
}
func NewSendContextImageRequest(to, imageId, imageURL, targetMessage string) *SendMessage {
	ib := &ImageBody{}

	if imageId != "" {
		ib.Id = &imageId
	} else {
		ib.Link = &imageURL
	}
	return &SendMessage{
		MessagingProduct: "whatsapp",
		To:               to,
		Type:             "image",
		ContextMessage:   &ContextMessage{Context: &Context{MessageId: targetMessage}},
		ImageMessage:     &ImageMessage{Image: ib},
	}
}
func (s *SendMessage) validateImageMessage() error {
	return nil
}
