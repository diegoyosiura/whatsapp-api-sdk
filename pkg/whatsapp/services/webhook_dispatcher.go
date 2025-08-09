package services

import "github.com/diegoyosiura/whatsapp-sdk-go/pkg/whatsapp/domain"

type WebhookHandler interface {
	OnMessage(m domain.InboundMessage)
	OnStatus(s domain.MessageStatus)
}

type WebhookDispatcher struct{ h WebhookHandler }

func NewWebhookDispatcher(h WebhookHandler) *WebhookDispatcher { return &WebhookDispatcher{h: h} }

func (d *WebhookDispatcher) Dispatch(e domain.WebhookEvent) {
	for _, entry := range e.Entry {
		for _, ch := range entry.Changes {
			if len(ch.Value.Messages) > 0 {
				for _, m := range ch.Value.Messages {
					d.h.OnMessage(m)
				}
			}
			if len(ch.Value.Statuses) > 0 {
				for _, s := range ch.Value.Statuses {
					d.h.OnStatus(s)
				}
			}
		}
	}
}
