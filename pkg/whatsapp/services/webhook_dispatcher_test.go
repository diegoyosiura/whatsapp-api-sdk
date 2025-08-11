package services_test

import (
	"testing"

	"github.com/diegoyosiura/whatsapp-sdk-go/pkg/whatsapp/domain"
	"github.com/diegoyosiura/whatsapp-sdk-go/pkg/whatsapp/services"
)

type fakeWebhookHandler struct {
	messages []domain.InboundMessage
	statuses []domain.MessageStatus
}

func (f *fakeWebhookHandler) OnMessage(m domain.InboundMessage) { f.messages = append(f.messages, m) }
func (f *fakeWebhookHandler) OnStatus(s domain.MessageStatus)   { f.statuses = append(f.statuses, s) }

func TestWebhookDispatcher_Dispatch(t *testing.T) {
	h := &fakeWebhookHandler{}
	dispatcher := services.NewWebhookDispatcher(h)

	event := domain.WebhookEvent{
		Entry: []domain.WebhookEntry{{
			Changes: []domain.WebhookChange{
				{Value: domain.WebhookValue{Messages: []domain.InboundMessage{{ID: "m1"}, {ID: "m2"}}}},
				{Value: domain.WebhookValue{Statuses: []domain.MessageStatus{{ID: "s1"}, {ID: "s2"}}}},
			},
		}},
	}

	dispatcher.Dispatch(event)

	if len(h.messages) != 2 {
		t.Fatalf("expected 2 messages, got %d", len(h.messages))
	}
	if len(h.statuses) != 2 {
		t.Fatalf("expected 2 statuses, got %d", len(h.statuses))
	}
	if h.messages[0].ID != "m1" || h.messages[1].ID != "m2" {
		t.Fatalf("unexpected messages: %+v", h.messages)
	}
	if h.statuses[0].ID != "s1" || h.statuses[1].ID != "s2" {
		t.Fatalf("unexpected statuses: %+v", h.statuses)
	}
}

func TestWebhookDispatcher_DispatchEmpty(t *testing.T) {
	h := &fakeWebhookHandler{}
	dispatcher := services.NewWebhookDispatcher(h)
	dispatcher.Dispatch(domain.WebhookEvent{})
	if len(h.messages) != 0 {
		t.Fatalf("expected no messages, got %d", len(h.messages))
	}
	if len(h.statuses) != 0 {
		t.Fatalf("expected no statuses, got %d", len(h.statuses))
	}
}
