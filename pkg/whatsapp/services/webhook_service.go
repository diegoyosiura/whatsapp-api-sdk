package services

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/diegoyosiura/whatsapp-sdk-go/pkg/whatsapp/domain"
	"github.com/diegoyosiura/whatsapp-sdk-go/pkg/whatsapp/ports"
)

// WebhookService encapsulates webhook validation and parsing logic.
// It is transport-agnostic and depends only on hexagonal ports.
type WebhookService struct {
	secrets ports.SecretsProvider
}

func NewWebhookService(secrets ports.SecretsProvider) *WebhookService {
	return &WebhookService{secrets: secrets}
}

// ValidateVerifyToken checks the provided token against the stored verify token.
// Use this in the GET verification handshake: the provider echoes back hub.challenge
// only if the verify token matches.
func (s *WebhookService) ValidateVerifyToken(ctx context.Context, provided string) error {
	if strings.TrimSpace(provided) == "" {
		return errors.New("verify token is empty")
	}
	want, err := s.secrets.Get(ctx, ports.VerifyTokenKey)
	if err != nil {
		return fmt.Errorf("failed to read verify token: %w", err)
	}
	if provided != want {
		return errors.New("verify token mismatch")
	}
	return nil
}

// VerifySignature validates the X-Hub-Signature-256 HMAC header against the raw body.
// Expected header format: "sha256=<hex_digest>".
func (s *WebhookService) VerifySignature(ctx context.Context, rawBody []byte, header string) error {
	if len(rawBody) == 0 {
		return errors.New("empty body")
	}
	if !strings.HasPrefix(header, "sha256=") || len(header) <= len("sha256=") {
		return errors.New("invalid signature header format")
	}

	secret, err := s.secrets.Get(ctx, ports.AppSecretKey)
	if err != nil {
		return fmt.Errorf("failed to read app secret: %w", err)
	}

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(rawBody)
	want := mac.Sum(nil)

	haveHex := header[len("sha256="):]
	have, err := hex.DecodeString(haveHex)
	if err != nil {
		return fmt.Errorf("invalid signature hex: %w", err)
	}

	if !hmac.Equal(want, have) {
		return errors.New("signature mismatch")
	}
	return nil
}

// ParseEvent decodes the JSON payload into the domain model.
func (s *WebhookService) ParseEvent(_ context.Context, rawBody []byte) (domain.WebhookEvent, error) {
	return domain.ParseWebhookEvent(rawBody)
}
