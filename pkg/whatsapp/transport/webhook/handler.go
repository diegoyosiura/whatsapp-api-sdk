package webhook

import (
	"io"
	"net/http"

	"github.com/diegoyosiura/whatsapp-sdk-go/pkg/whatsapp/services"
)

// Handler is an HTTP adapter for WhatsApp webhooks.
// It uses WebhookService for validation/parsing and optionally dispatches
// events to a WebhookDispatcher.
type Handler struct {
	Service    *services.WebhookService
	Dispatcher *services.WebhookDispatcher
}

func NewHandler(svc *services.WebhookService, dispatcher *services.WebhookDispatcher) *Handler {
	return &Handler{Service: svc, Dispatcher: dispatcher}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	switch r.Method {
	case http.MethodGet:
		mode := r.URL.Query().Get("hub.mode")
		token := r.URL.Query().Get("hub.verify_token")
		challenge := r.URL.Query().Get("hub.challenge")
		if mode != "subscribe" {
			http.Error(w, "invalid mode", http.StatusBadRequest)
			return
		}
		if err := h.Service.ValidateVerifyToken(ctx, token); err != nil {
			http.Error(w, "verify token mismatch", http.StatusForbidden)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, challenge)
		return
	case http.MethodPost:
		raw, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "cannot read body", http.StatusBadRequest)
			return
		}
		sig := r.Header.Get("X-Hub-Signature-256")
		if err := h.Service.VerifySignature(ctx, raw, sig); err != nil {
			http.Error(w, "invalid signature", http.StatusForbidden)
			return
		}
		event, err := h.Service.ParseEvent(ctx, raw)
		if err != nil {
			http.Error(w, "invalid payload", http.StatusBadRequest)
			return
		}
		if h.Dispatcher != nil {
			h.Dispatcher.Dispatch(event, r.Header)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, "ok")
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}
