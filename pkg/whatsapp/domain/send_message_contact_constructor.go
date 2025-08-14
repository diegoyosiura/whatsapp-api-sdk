package domain

import "github.com/diegoyosiura/whatsapp-sdk-go/pkg/errorsx"

type ContactMessage struct {
	Contacts []*Contact `json:"contacts"`
}

type Contact struct {
	Birthday  string            `json:"birthday"`
	Name      *ContactName      `json:"name"`
	Org       *ContactOrg       `json:"org"`
	Addresses []*ContactAddress `json:"addresses"`
	Emails    []*ContactEmail   `json:"emails"`
	Phones    []*ContactPhone   `json:"phones"`
	Urls      []*ContactURL     `json:"urls"`
}

type ContactAddress struct {
	Street      string `json:"street"`
	City        string `json:"city"`
	State       string `json:"state"`
	Zip         string `json:"zip"`
	Country     string `json:"country"`
	CountryCode string `json:"country_code"`
	Type        string `json:"type"`
}
type ContactEmail struct {
	Email string `json:"email"`
	Type  string `json:"type"`
}
type ContactName struct {
	FormattedName string `json:"formatted_name"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	MiddleName    string `json:"middle_name"`
	Suffix        string `json:"suffix"`
	Prefix        string `json:"prefix"`
}

type ContactOrg struct {
	Company    string `json:"company"`
	Department string `json:"department"`
	Title      string `json:"title"`
}
type ContactPhone struct {
	Phone string `json:"phone"`
	WaId  string `json:"wa_id"`
	Type  string `json:"type"`
}
type ContactURL struct {
	Url  string `json:"url"`
	Type string `json:"type"`
}

func NewSendContactRequest(to string, contactList []*Contact) *SendMessage {
	return &SendMessage{
		MessagingProduct: "whatsapp",
		To:               to,
		Type:             "contacts",
		ContactMessage:   &ContactMessage{Contacts: contactList},
	}
}
func NewSendContextContactRequest(to, targetMessage string, contactList []*Contact) *SendMessage {
	return &SendMessage{
		MessagingProduct: "whatsapp",
		To:               to,
		Type:             "contacts",
		ContextMessage:   &ContextMessage{Context: &Context{MessageId: targetMessage}},
		ContactMessage:   &ContactMessage{Contacts: contactList},
	}
}

func (s *SendMessage) validateContactMessage() error {
	if s.Type != "contacts" {
		return &errorsx.ValidationError{Field: "Type", Reason: "type must be contacts", Op: "validateContactMessage"}
	}
	if s.ContactMessage == nil {
		return &errorsx.ValidationError{Field: "ContactMessage", Reason: "contact is nil", Op: "validateContactMessage"}
	}
	if len(s.ContactMessage.Contacts) == 0 {
		return &errorsx.ValidationError{Field: "Contacts", Reason: "must have at least one contact", Op: "validateContactMessage"}
	}

	for _, contact := range s.ContactMessage.Contacts {
		if contact == nil {
			return &errorsx.ValidationError{Field: "Contacts[]", Reason: "nil contact found", Op: "validateContactMessage"}
		}
		if contact.Name == nil {
			return &errorsx.ValidationError{Field: "Contacts[].Name", Reason: "nil contact name found", Op: "validateContactMessage"}
		}

		if contact.Name.FirstName == "" {
			return &errorsx.ValidationError{Field: "Contacts[].Name.FirstName", Reason: "contact.Name.FirstName is empty", Op: "validateContactMessage"}
		}
	}
	return nil
}
