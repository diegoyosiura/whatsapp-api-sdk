package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/diegoyosiura/whatsapp-sdk-go/pkg/whatsapp/ports"
	"github.com/diegoyosiura/whatsapp-sdk-go/pkg/whatsapp/services"
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

func main() {
	secrets := envSecrets{}
	svc := services.NewWebhookService(secrets)

	http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		switch r.Method {
		case http.MethodGet:
			// Verification handshake
			mode := r.URL.Query().Get("hub.mode")
			token := r.URL.Query().Get("hub.verify_token")
			challenge := r.URL.Query().Get("hub.challenge")

			if mode != "subscribe" {
				w.WriteHeader(http.StatusBadRequest)
				_, _ = w.Write([]byte("invalid mode"))
				return
			}
			if err := svc.ValidateVerifyToken(ctx, token); err != nil {
				w.WriteHeader(http.StatusForbidden)
				_, _ = w.Write([]byte("verify token mismatch"))
				return
			}
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(challenge))
			return

		case http.MethodPost:
			// Signature verification + parse
			raw, err := io.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				_, _ = w.Write([]byte("cannot read body"))
				return
			}
			sig := r.Header.Get("X-Hub-Signature-256")
			if err := svc.VerifySignature(ctx, raw, sig); err != nil {
				w.WriteHeader(http.StatusForbidden)
				_, _ = w.Write([]byte("invalid signature"))
				return
			}

			event, err := svc.ParseEvent(ctx, raw)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				_, _ = w.Write([]byte("invalid payload"))
				return
			}

			// Business hook: log basic info. Replace with your dispatcher.
			log.Printf("incoming webhook: object=%s entries=%d", event.Object, len(event.Entry))

			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("ok"))
			return

		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
	})

	addr := ":8080"
	log.Printf("listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
