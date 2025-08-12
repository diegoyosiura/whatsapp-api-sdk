package services

import (
	"net/http"

	"github.com/diegoyosiura/whatsapp-sdk-go/pkg/whatsapp/domain"
)

type WebhookHandler interface {
	Always(e domain.WebhookEvent, h http.Header)
	OnMessage(m domain.InboundMessage, e domain.WebhookEvent, h http.Header)
	OnStatus(s domain.MessageStatus, e domain.WebhookEvent, h http.Header)
}

type WebhookDispatcher struct{ h WebhookHandler }

func NewWebhookDispatcher(h WebhookHandler) *WebhookDispatcher { return &WebhookDispatcher{h: h} }

func (d *WebhookDispatcher) Dispatch(e domain.WebhookEvent, h http.Header) {
	d.h.Always(e, h)
	for _, entry := range e.Entry {
		for _, ch := range entry.Changes {
			if len(ch.Value.Messages) > 0 {
				for _, m := range ch.Value.Messages {
					d.h.OnMessage(m, e, h)
				}
			}
			if len(ch.Value.Statuses) > 0 {
				for _, s := range ch.Value.Statuses {
					d.h.OnStatus(s, e, h)
				}
			}
		}
	}
}
