package domain

import "testing"

func newValidLocationMessage() *SendMessage {
	return NewSendLocationRequest("+123456789", "1", "2", "Loc", "Addr")
}

func TestNewSendLocationRequest(t *testing.T) {
	msg := newValidLocationMessage()
	if msg.Type != "location" {
		t.Fatalf("unexpected type %s", msg.Type)
	}
	if msg.LocationMessage == nil || msg.LocationMessage.Location == nil {
		t.Fatalf("location not set")
	}
	if err := msg.Validate(); err != nil {
		t.Fatalf("Validate error: %v", err)
	}
}

func TestNewSendContextLocationRequest(t *testing.T) {
	msg := NewSendContextLocationRequest("+123456789", "1", "2", "Loc", "Addr", "m1")
	if msg.ContextMessage == nil || msg.ContextMessage.Context.MessageId != "m1" {
		t.Fatalf("context not set")
	}
}

func TestValidateLocationMessageErrors(t *testing.T) {
	tests := []struct {
		name   string
		modify func(*SendMessage)
	}{
		{"wrong type", func(m *SendMessage) { m.Type = "text" }},
		{"nil body", func(m *SendMessage) { m.LocationMessage = nil }},
		{"nil location", func(m *SendMessage) { m.LocationMessage.Location = nil }},
		{"empty name", func(m *SendMessage) { m.LocationMessage.Location.Name = "" }},
		{"empty latitude", func(m *SendMessage) { m.LocationMessage.Location.Latitude = "" }},
		{"empty longitude", func(m *SendMessage) { m.LocationMessage.Location.Longitude = "" }},
		{"empty address", func(m *SendMessage) { m.LocationMessage.Location.Address = "" }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := newValidLocationMessage()
			tt.modify(msg)
			if err := msg.validateLocationMessage(); err == nil {
				t.Fatalf("expected error")
			}
		})
	}
}
