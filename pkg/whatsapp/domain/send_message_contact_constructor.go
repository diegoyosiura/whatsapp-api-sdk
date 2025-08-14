package domain

type ContactMessage struct {
	Contacts []*Contact `json:"contacts"`
}

type Contact struct {
	Addresses []struct {
		Street      string `json:"street"`
		City        string `json:"city"`
		State       string `json:"state"`
		Zip         string `json:"zip"`
		Country     string `json:"country"`
		CountryCode string `json:"country_code"`
		Type        string `json:"type"`
	} `json:"addresses"`
	Birthday string `json:"birthday"`
	Emails   []struct {
		Email string `json:"email"`
		Type  string `json:"type"`
	} `json:"emails"`
	Name struct {
		FormattedName string `json:"formatted_name"`
		FirstName     string `json:"first_name"`
		LastName      string `json:"last_name"`
		MiddleName    string `json:"middle_name"`
		Suffix        string `json:"suffix"`
		Prefix        string `json:"prefix"`
	} `json:"name"`
	Org struct {
		Company    string `json:"company"`
		Department string `json:"department"`
		Title      string `json:"title"`
	} `json:"org"`
	Phones []struct {
		Phone string `json:"phone"`
		WaId  string `json:"wa_id"`
		Type  string `json:"type"`
	} `json:"phones"`
	Urls []struct {
		Url  string `json:"url"`
		Type string `json:"type"`
	} `json:"urls"`
}

func NewSendContactRequest(to, body string) *SendMessage {
	return &SendMessage{
		MessagingProduct: "whatsapp",
		To:               to,
		Type:             "text",
		ContactMessage:   &ContactMessage{Contacts: make([]*Contact, 0)},
	}
}
func NewSendContextContactRequest(to, body, targetMessage string) *SendMessage {
	return &SendMessage{
		MessagingProduct: "whatsapp",
		To:               to,
		Type:             "text",
		ContextMessage:   &ContextMessage{Context: &Context{MessageId: targetMessage}},
		ContactMessage:   &ContactMessage{Contacts: make([]*Contact, 0)},
	}
}

func (s *SendMessage) validateContactMessage() error {
	return nil
}
