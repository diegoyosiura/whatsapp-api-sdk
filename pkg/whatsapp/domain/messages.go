package domain

// MessageSendResponse represents a successful response from the WhatsApp
// messages endpoint when sending a message. The exact fields mirror the public
// API but are kept minimal here; they can be expanded as we implement more
// message types.
//
// Reference: WhatsApp Cloud API send message response shape.
type MessageSendResponse struct {
	MessagingProduct string `json:"messaging_product"`
	Contacts         []struct {
		Input string `json:"input"`
		WaID  string `json:"wa_id"`
	} `json:"contacts"`
	Messages []struct {
		ID string `json:"id"`
	} `json:"messages"`
}
