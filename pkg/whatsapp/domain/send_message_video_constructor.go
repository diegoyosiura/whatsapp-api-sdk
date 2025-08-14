package domain

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
	return nil
}
