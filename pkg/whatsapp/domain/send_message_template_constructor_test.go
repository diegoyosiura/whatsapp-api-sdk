package domain

import "testing"

func newValidTemplateMessage() *SendMessage {
	text := "hi"
	comp := &TemplateComponent{Type: "body", Parameters: []*TemplateParameter{{Type: "text", Text: &text}}}
	return NewSendTemplateRequest("+123456789", "tmpl", "en", []*TemplateComponent{comp})
}

func TestNewSendTemplateRequest(t *testing.T) {
	msg := newValidTemplateMessage()
	if msg.Type != "template" {
		t.Fatalf("unexpected type %s", msg.Type)
	}
	if msg.TemplateMessage == nil || msg.TemplateMessage.Template == nil {
		t.Fatalf("template not set")
	}
	if err := msg.Validate(); err != nil {
		t.Fatalf("Validate error: %v", err)
	}
}

func TestValidateTemplateMessageErrors(t *testing.T) {
	tests := []struct {
		name   string
		modify func(*SendMessage)
	}{
		{"nil body", func(m *SendMessage) { m.TemplateMessage = nil }},
		{"nil template", func(m *SendMessage) { m.TemplateMessage.Template = nil }},
		{"nil language", func(m *SendMessage) { m.TemplateMessage.Template.Language = nil }},
		{"empty components", func(m *SendMessage) { m.TemplateMessage.Template.Components = nil }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := newValidTemplateMessage()
			tt.modify(msg)
			if err := msg.validateTemplateMessage(); err == nil {
				t.Fatalf("expected error")
			}
		})
	}
}
