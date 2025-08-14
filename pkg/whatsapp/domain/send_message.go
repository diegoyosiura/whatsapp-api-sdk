package domain

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/diegoyosiura/whatsapp-sdk-go/pkg/errorsx"
	"github.com/diegoyosiura/whatsapp-sdk-go/pkg/utils"
)

// SendMessage is a minimal shape for sending a text message.
// The complete shape will live in the domain layer; here we only need a
// transport-level representation for the request body.
type SendMessage struct {
	MessagingProduct string  `json:"messaging_product"`
	RecipientType    *string `json:"recipient_type"`
	To               string  `json:"to"`
	Type             string  `json:"type"`

	*ContextMessage

	*TextMessage
	*ImageMessage
	*AudioMessage
	*ReactionMessage
	*DocumentMessage
	*StickerMessage
	*VideoMessage

	*ContactMessage
	*LocationMessage
	*TemplateMessage
}

func (s *SendMessage) Buffer() (*bytes.Buffer, error) {
	buf := &bytes.Buffer{}
	if err := json.NewEncoder(buf).Encode(s); err != nil {
		return nil, fmt.Errorf("encode request: %w", err)
	}

	return buf, nil
}

func (s *SendMessage) Validate() error {
	if s.To == "" {
		return &errorsx.ValidationError{Op: "SendMessage", Field: "to", Reason: "empty"}
	}
	if !utils.IsE164(s.To) {
		return &errorsx.ValidationError{Op: "SendMessage", Field: "to", Reason: "must be E.164 like +5511999999999"}
	}

	if s.ContextMessage != nil {
		if s.ContextMessage.Context.MessageId == "" {
			return &errorsx.ValidationError{Op: "SendMessage", Field: "Context.MessageId", Reason: "should not be empty"}
		}
	}

	switch s.Type {
	case "text":
		s.TextMessage.Text.PreviewURL = utils.HasURL(s.TextMessage.Text.Body)
		if err := s.validateTextMessage(); err != nil {
			return err
		}
		break
	case "image":
		if err := s.validateImageMessage(); err != nil {
			return err
		}
		break
	case "audio":
		if err := s.validateAudioMessage(); err != nil {
			return err
		}
		break
	case "reaction":
		if err := s.validateReplyMessage(); err != nil {
			return err
		}
		break
	case "document":
		if err := s.validateDocumentMessage(); err != nil {
			return err
		}
		break
	case "sticker":
		if err := s.validateStickerMessage(); err != nil {
			return err
		}
		break
	case "video":
		if err := s.validateVideoMessage(); err != nil {
			return err
		}
		break
	case "contacts":
		if err := s.validateContactMessage(); err != nil {
			return err
		}
		break
	case "location":
		if err := s.validateLocationMessage(); err != nil {
			return err
		}
		break
	case "template":
		if err := s.validateTemplateMessage(); err != nil {
			return err
		}
		break
	}
	return nil
}
