package domain

import "testing"

func newValidContactMessage() *SendMessage {
	return NewSendContactRequest("+123456789", []*Contact{{Name: &ContactName{FirstName: "John"}}})
}

func TestNewSendContactRequest(t *testing.T) {
	msg := newValidContactMessage()
	if msg.Type != "contacts" {
		t.Fatalf("unexpected type %s", msg.Type)
	}
	if msg.ContactMessage == nil || len(msg.ContactMessage.Contacts) != 1 {
		t.Fatalf("contacts not set")
	}
	if msg.ContactMessage.Contacts[0].Name.FirstName != "John" {
		t.Fatalf("first name not set")
	}
	if err := msg.Validate(); err != nil {
		t.Fatalf("Validate error: %v", err)
	}
}

func TestNewSendContextContactRequest(t *testing.T) {
	msg := NewSendContextContactRequest("+123456789", "m1", []*Contact{{Name: &ContactName{FirstName: "John"}}})
	if msg.ContextMessage == nil || msg.ContextMessage.Context.MessageId != "m1" {
		t.Fatalf("context not set")
	}
}

func TestValidateContactMessageErrors(t *testing.T) {
	tests := []struct {
		name   string
		modify func(*SendMessage)
	}{
		{"wrong type", func(m *SendMessage) { m.Type = "text" }},
		{"nil body", func(m *SendMessage) { m.ContactMessage = nil }},
		{"empty contacts", func(m *SendMessage) { m.ContactMessage.Contacts = nil }},
		{"nil contact", func(m *SendMessage) { m.ContactMessage.Contacts = []*Contact{nil} }},
		{"nil name", func(m *SendMessage) { m.ContactMessage.Contacts = []*Contact{{Name: nil}} }},
		{"empty first name", func(m *SendMessage) { m.ContactMessage.Contacts = []*Contact{{Name: &ContactName{}}} }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := newValidContactMessage()
			tt.modify(msg)
			if err := msg.validateContactMessage(); err == nil {
				t.Fatalf("expected error")
			}
		})
	}
}
