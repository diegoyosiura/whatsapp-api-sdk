package domain

import "github.com/diegoyosiura/whatsapp-sdk-go/pkg/errorsx"

type VideoMessage struct {
	Video *VideoBody `json:"video"`
}
type VideoBody struct {
	Id      *string `json:"id"`
	Link    *string `json:"link"`
	Caption string  `json:"caption"`
}

func NewSendVideoRequest(to, videoID, videoLink, caption string) *SendMessage {
	rt := "individual"

	vid := &VideoBody{
		Caption: caption,
	}

	if videoID != "" {
		vid.Id = &videoID
	} else {
		vid.Link = &videoLink
	}

	return &SendMessage{
		MessagingProduct: "whatsapp",
		RecipientType:    &rt,
		To:               to,
		Type:             "video",
		VideoMessage:     &VideoMessage{Video: vid},
	}
}
func NewSendContextVideoRequest(to, videoID, videoLink, caption, targetMessage string) *SendMessage {
	vid := &VideoBody{
		Caption: caption,
	}

	if videoID != "" {
		vid.Id = &videoID
	} else {
		vid.Link = &videoLink
	}

	return &SendMessage{
		MessagingProduct: "whatsapp",
		To:               to,
		Type:             "video",
		ContextMessage:   &ContextMessage{Context: &Context{MessageId: targetMessage}},
		VideoMessage:     &VideoMessage{Video: vid},
	}
}
func (s *SendMessage) validateVideoMessage() error {
	if s.Type != "video" {
		return &errorsx.ValidationError{Field: "Type", Reason: "type must be video", Op: "validateVideoMessage"}
	}

	if s.VideoMessage == nil {
		return &errorsx.ValidationError{Field: "VideoMessage", Reason: "video is nil", Op: "validateVideoMessage"}
	}

	if s.VideoMessage.Video == nil {
		return &errorsx.ValidationError{Field: "Video", Reason: "video is nil", Op: "validateVideoMessage"}
	}

	if *s.VideoMessage.Video.Id == "" && *s.VideoMessage.Video.Link == "" {
		return &errorsx.ValidationError{Field: "Video", Reason: "the reference must be a link or a Id, nothing received", Op: "validateVideoMessage"}
	}

	if *s.VideoMessage.Video.Id != "" && *s.VideoMessage.Video.Link != "" {
		return &errorsx.ValidationError{Field: "Video", Reason: "the reference must be a link or a Id, both received", Op: "validateVideoMessage"}
	}
	return nil
}
