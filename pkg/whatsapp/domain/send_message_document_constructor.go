package domain

type DocumentMessage struct {
	Document *DocumentBody `json:"document"`
}
type DocumentBody struct {
	Id       *string `json:"id"`
	Link     *string `json:"link"`
	Caption  string  `json:"caption"`
	Filename string  `json:"filename"`
}

func NewSendDocumentRequest(to, documentID, documentLink, caption, fileName string) *SendMessage {
	rt := "individual"

	doc := &DocumentBody{
		Caption:  caption,
		Filename: fileName,
	}

	if documentID != "" {
		doc.Id = &documentID
	} else {
		doc.Link = &documentLink
	}

	return &SendMessage{
		MessagingProduct: "whatsapp",
		RecipientType:    &rt,
		To:               to,
		Type:             "document",
		DocumentMessage:  &DocumentMessage{Document: doc},
	}
}
func NewSendContextDocumentRequest(to, documentID, documentLink, caption, fileName, targetMessage string) *SendMessage {
	doc := &DocumentBody{
		Caption:  caption,
		Filename: fileName,
	}

	if documentID != "" {
		doc.Id = &documentID
	} else {
		doc.Link = &documentLink
	}

	return &SendMessage{
		MessagingProduct: "whatsapp",
		To:               to,
		Type:             "document",
		ContextMessage:   &ContextMessage{Context: &Context{MessageId: targetMessage}},
		DocumentMessage:  &DocumentMessage{Document: doc},
	}
}
func (s *SendMessage) validateDocumentMessage() error {
	return nil
}
