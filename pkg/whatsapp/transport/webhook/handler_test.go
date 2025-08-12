package webhook_test

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"testing"

	intfakes "github.com/diegoyosiura/whatsapp-sdk-go/internal/testutils/whatsapp/ports"
	"github.com/diegoyosiura/whatsapp-sdk-go/pkg/whatsapp/domain"
	"github.com/diegoyosiura/whatsapp-sdk-go/pkg/whatsapp/ports"
	"github.com/diegoyosiura/whatsapp-sdk-go/pkg/whatsapp/services"
	"github.com/diegoyosiura/whatsapp-sdk-go/pkg/whatsapp/transport/webhook"
)

type fakeHandler struct{ messages []domain.InboundMessage }

func (f *fakeHandler) OnMessage(m domain.InboundMessage, e domain.WebhookEvent, h http.Header) {
	f.messages = append(f.messages, m)
}
func (f *fakeHandler) OnStatus(s domain.MessageStatus, e domain.WebhookEvent, h http.Header) {}

func TestHandler_GetVerify(t *testing.T) {
	fp := &intfakes.FakeSecretsProvider{Secrets: map[ports.SecretKey]string{
		ports.VerifyTokenKey: "tkn",
	}}
	svc := services.NewWebhookService(fp)
	h := webhook.NewHandler(svc, nil)

	req := httptest.NewRequest(http.MethodGet, "/webhook?hub.mode=subscribe&hub.verify_token=tkn&hub.challenge=abc", nil)
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK || rr.Body.String() != "abc" {
		t.Fatalf("unexpected response: code=%d body=%s", rr.Code, rr.Body.String())
	}
}

func TestHandler_PostDispatch(t *testing.T) {
	payload := []byte(`{"object":"whatsapp_business_account","entry":[{"id":"x","changes":[{"value":{"messages":[{"id":"m1"}]}}]}]}`)
	appSecret := "secret"
	mac := hmac.New(sha256.New, []byte(appSecret))
	mac.Write(payload)
	sig := "sha256=" + hex.EncodeToString(mac.Sum(nil))

	fp := &intfakes.FakeSecretsProvider{Secrets: map[ports.SecretKey]string{
		ports.AppSecretKey: appSecret,
	}}
	svc := services.NewWebhookService(fp)
	fh := &fakeHandler{}
	dispatcher := services.NewWebhookDispatcher(fh)
	h := webhook.NewHandler(svc, dispatcher)

	req := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewReader(payload))
	req.Header.Set("X-Hub-Signature-256", sig)
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d", rr.Code)
	}
	if len(fh.messages) != 1 || fh.messages[0].ID != "m1" {
		t.Fatalf("dispatcher not invoked: %+v", fh.messages)
	}
}
