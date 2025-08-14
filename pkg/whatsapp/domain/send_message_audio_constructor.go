package domain

import (
	"github.com/diegoyosiura/whatsapp-sdk-go/pkg/errorsx"
)

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
	if s.Type != "audio" {
		return &errorsx.ValidationError{Field: "Type", Reason: "type must be audio", Op: "validateAudioMessage"}
	}

	if s.AudioMessage == nil {
		return &errorsx.ValidationError{Field: "AudioMessage", Reason: "audio is nil", Op: "validateAudioMessage"}
	}

	if s.AudioMessage.Audio == nil {
		return &errorsx.ValidationError{Field: "Audio", Reason: "audio is nil", Op: "validateAudioMessage"}
	}

	if *s.AudioMessage.Audio.Id == "" && *s.AudioMessage.Audio.Link == "" {
		return &errorsx.ValidationError{Field: "Audio", Reason: "the reference must be a link or a Id, nothing received", Op: "validateAudioMessage"}
	}

	if *s.AudioMessage.Audio.Id != "" && *s.AudioMessage.Audio.Link != "" {
		return &errorsx.ValidationError{Field: "Audio", Reason: "the reference must be a link or a Id, both received", Op: "validateAudioMessage"}
	}

	return nil
}
