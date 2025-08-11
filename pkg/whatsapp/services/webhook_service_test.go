package services_test

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"testing"

	intfakes "github.com/diegoyosiura/whatsapp-sdk-go/internal/testutils/whatsapp/ports"
	"github.com/diegoyosiura/whatsapp-sdk-go/pkg/whatsapp/ports"
	"github.com/diegoyosiura/whatsapp-sdk-go/pkg/whatsapp/services"
)

func TestWebhookService_ValidateVerifyToken_OK(t *testing.T) {
	payload := []byte(`{"ok":true}`)
	appSecret := "my-app-secret"

	// Gera HMAC sha256 do payload para montar o header correto
	mac := hmac.New(sha256.New, []byte(appSecret))
	mac.Write(payload)
	header := "sha256=" + hex.EncodeToString(mac.Sum(nil))

	fp := &intfakes.FakeSecretsProvider{
		Secrets: map[ports.SecretKey]string{
			ports.AppSecretKey: appSecret,
		},
	}
	svc := services.NewWebhookService(fp)

	if err := svc.VerifySignature(context.Background(), payload, header); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestWebhookService_ValidateVerifyToken_Mismatch(t *testing.T) {
	fp := &intfakes.FakeSecretsProvider{
		Secrets: map[ports.SecretKey]string{
			ports.VerifyTokenKey: "expected",
		},
	}
	svc := services.NewWebhookService(fp)

	if err := svc.ValidateVerifyToken(context.Background(), "wrong"); err == nil {
		t.Fatalf("expected error on mismatch, got nil")
	}
}

func TestWebhookService_VerifySignature_OK(t *testing.T) {
	payload := []byte(`{"ok":true}`)
	appSecret := "my-app-secret"

	// calc HMAC sha256 do payload para formar o header esperado
	mac := hmac.New(sha256.New, []byte(appSecret))
	mac.Write(payload)
	sig := hex.EncodeToString(mac.Sum(nil))
	header := "sha256=" + sig

	fp := &intfakes.FakeSecretsProvider{
		Secrets: map[ports.SecretKey]string{
			ports.AppSecretKey: appSecret,
		},
	}
	svc := services.NewWebhookService(fp)

	if err := svc.VerifySignature(context.Background(), payload, header); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestWebhookService_VerifySignature_Bad(t *testing.T) {
	appSecret := "app-secret-123"
	payload := []byte(`{"x":1}`)
	badSig := "sha256=00deadbeef"

	fp := &intfakes.FakeSecretsProvider{
		Secrets: map[ports.SecretKey]string{
			ports.AppSecretKey: appSecret,
		},
	}
	svc := services.NewWebhookService(fp)

	if err := svc.VerifySignature(context.Background(), payload, badSig); err == nil {
		t.Fatalf("expected error for bad signature, got nil")
	}
}

func TestWebhookService_ParseEvent_OK(t *testing.T) {
	fp := &intfakes.FakeSecretsProvider{}
	svc := services.NewWebhookService(fp)

	raw := []byte(`{"object":"whatsapp_business_account","entry":[{"id":"x","changes":[]}]}`)
	ev, err := svc.ParseEvent(context.Background(), raw)
	if err != nil {
		t.Fatalf("ParseEvent error: %v", err)
	}
	if ev.Object != "whatsapp_business_account" {
		t.Fatalf("unexpected object: %s", ev.Object)
	}
}

func TestWebhookService_ValidateVerifyToken_Success(t *testing.T) {
	fp := &intfakes.FakeSecretsProvider{Secrets: map[ports.SecretKey]string{
		ports.VerifyTokenKey: "secret",
	}}
	svc := services.NewWebhookService(fp)
	if err := svc.ValidateVerifyToken(context.Background(), "secret"); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestWebhookService_ValidateVerifyToken_Empty(t *testing.T) {
	fp := &intfakes.FakeSecretsProvider{}
	svc := services.NewWebhookService(fp)
	if err := svc.ValidateVerifyToken(context.Background(), ""); err == nil {
		t.Fatalf("expected error for empty token, got nil")
	}
}

func TestWebhookService_VerifySignature_InvalidHeader(t *testing.T) {
	fp := &intfakes.FakeSecretsProvider{Secrets: map[ports.SecretKey]string{
		ports.AppSecretKey: "secret",
	}}
	svc := services.NewWebhookService(fp)
	if err := svc.VerifySignature(context.Background(), []byte("body"), "bad"); err == nil {
		t.Fatalf("expected error for invalid header, got nil")
	}
}

func TestWebhookService_VerifySignature_InvalidHex(t *testing.T) {
	fp := &intfakes.FakeSecretsProvider{Secrets: map[ports.SecretKey]string{
		ports.AppSecretKey: "secret",
	}}
	svc := services.NewWebhookService(fp)
	if err := svc.VerifySignature(context.Background(), []byte("body"), "sha256=zzzz"); err == nil {
		t.Fatalf("expected error for invalid hex, got nil")
	}
}

func TestWebhookService_ParseEvent_BadJSON(t *testing.T) {
	fp := &intfakes.FakeSecretsProvider{}
	svc := services.NewWebhookService(fp)
	if _, err := svc.ParseEvent(context.Background(), []byte("not-json")); err == nil {
		t.Fatalf("expected error for bad json, got nil")
	}
}
