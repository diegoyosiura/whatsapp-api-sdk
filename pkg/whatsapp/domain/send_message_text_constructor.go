package domain

import "github.com/diegoyosiura/whatsapp-sdk-go/pkg/errorsx"

type TextMessage struct {
	Text *TextBody `json:"text"`
}
type TextBody struct {
	PreviewURL bool   `json:"preview_url"`
	Body       string `json:"body"`
}

func NewSendTestMessage(to, body string) *SendMessage {
	return &SendMessage{
		MessagingProduct: "whatsapp",
		To:               to,
		Type:             "text",
		TextMessage:      &TextMessage{Text: &TextBody{Body: body}},
	}
}
func NewSendContextTextRequest(to, body, targetMessage string) *SendMessage {
	return &SendMessage{
		MessagingProduct: "whatsapp",
		To:               to,
		Type:             "text",
		ContextMessage:   &ContextMessage{Context: &Context{MessageId: targetMessage}},
		TextMessage:      &TextMessage{Text: &TextBody{Body: body}},
	}
}

func (s *SendMessage) validateTextMessage() error {
	if s.TextMessage.Text.Body == "" {
		return &errorsx.ValidationError{Op: "SendMessage", Field: "body", Reason: "empty"}
	}

	return nil
}
