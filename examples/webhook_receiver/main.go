package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/diegoyosiura/whatsapp-sdk-go/pkg/whatsapp/domain"
	"github.com/diegoyosiura/whatsapp-sdk-go/pkg/whatsapp/ports"
	"github.com/diegoyosiura/whatsapp-sdk-go/pkg/whatsapp/services"
	"github.com/diegoyosiura/whatsapp-sdk-go/pkg/whatsapp/transport/webhook"
)

// Minimal SecretsProvider backed by environment variables for the example only.
type envSecrets struct{}

func (s envSecrets) Get(ctx context.Context, key ports.SecretKey) (string, error) {
	v := os.Getenv(string(key))
	if v == "" {
		return "", fmt.Errorf("secret %s not set", key)
	}
	return v, nil
}

// logHandler logs webhook events; replace with your business logic.
type logHandler struct{}

func (logHandler) OnMessage(m domain.InboundMessage) { log.Printf("incoming message id=%s", m.ID) }
func (logHandler) OnStatus(s domain.MessageStatus) {
	log.Printf("status id=%s status=%s", s.ID, s.Status)
}

func main() {
	secrets := envSecrets{}
	svc := services.NewWebhookService(secrets)
	dispatcher := services.NewWebhookDispatcher(logHandler{})
	http.Handle("/webhook", webhook.NewHandler(svc, dispatcher))

	addr := ":8080"
	log.Printf("listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
