package domain

import "github.com/diegoyosiura/whatsapp-sdk-go/pkg/errorsx"

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
	if s.Type != "image" {
		return &errorsx.ValidationError{Field: "Type", Reason: "type must be image", Op: "validateImageMessage"}
	}

	if s.ImageMessage == nil {
		return &errorsx.ValidationError{Field: "ImageMessage", Reason: "image is nil", Op: "validateImageMessage"}
	}

	if s.ImageMessage.Image == nil {
		return &errorsx.ValidationError{Field: "Image", Reason: "image is nil", Op: "validateImageMessage"}
	}

	if *s.ImageMessage.Image.Id == "" && *s.ImageMessage.Image.Link == "" {
		return &errorsx.ValidationError{Field: "Image", Reason: "the reference must be a link or a Id, nothing received", Op: "validateImageMessage"}
	}

	if *s.ImageMessage.Image.Id != "" && *s.ImageMessage.Image.Link != "" {
		return &errorsx.ValidationError{Field: "Image", Reason: "the reference must be a link or a Id, both received", Op: "validateImageMessage"}
	}
	return nil
}
