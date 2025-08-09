package services_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/diegoyosiura/whatsapp-sdk-go/pkg/errorsx"
	"github.com/diegoyosiura/whatsapp-sdk-go/pkg/whatsapp"
	"github.com/diegoyosiura/whatsapp-sdk-go/pkg/whatsapp/services"
)

// fakeTokenProvider is a lightweight test double used locally in this package
// to avoid importing internal testutils from another module path.
type fakeTokenProvider struct{ token string }

func (f fakeTokenProvider) Token(ctx context.Context) (string, error) { return f.token, nil }
func (f fakeTokenProvider) Refresh(ctx context.Context) error         { return nil }

// newTestClient constructs a whatsapp.Client pointing to the provided baseURL
// and using the internal httpx default (by passing nil HTTPDoer). It sets a
// small timeout and retry budget for tests.
func newTestClient(t *testing.T, baseURL string) *whatsapp.Client {
	t.Helper()
	opts := whatsapp.Options{
		Version:         "v20.0",
		WABAID:          "waba-test",
		PhoneNumberID:   "1234567890",
		HTTPDoer:        nil, // use default httpx with retry
		TokenProvider:   fakeTokenProvider{token: "test-token"},
		SecretsProvider: nil,
		BaseURL:         baseURL,
		Timeout:         2 * time.Second,
		RetryMax:        2,
		UserAgent:       "test-suite",
	}
	c, err := whatsapp.NewClient(opts)
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	return c
}

func TestMessagesService_SendText_Success(t *testing.T) {
	// Load success fixture
	fixturePath := filepath.Join("..", "..", "..", "testdata", "send_text_success.json")
	b, err := os.ReadFile(fixturePath)
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}

	// Prepare server that returns 200 with the fixture
	var capturedPath string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(b)
	}))
	defer ts.Close()

	c := newTestClient(t, ts.URL)
	svc := services.NewMessagesService(c)

	resp, err := svc.SendText(context.Background(), "+5511999999999", "hello")
	if err != nil {
		t.Fatalf("SendText: unexpected error: %v", err)
	}
	if resp == nil || len(resp.Messages) == 0 {
		t.Fatalf("expected message id in response, got %+v", resp)
	}
	// Ensure the path matches /v20.0/{PN-ID}/messages
	expectedPath := "/v20.0/1234567890/messages"
	if capturedPath != expectedPath {
		t.Fatalf("unexpected path: want %s got %s", expectedPath, capturedPath)
	}
}

func TestMessagesService_SendText_Validation(t *testing.T) {
	c := newTestClient(t, "http://invalid.local")
	svc := services.NewMessagesService(c)
	if _, err := svc.SendText(context.Background(), "", "body"); err == nil {
		t.Fatalf("expected validation error for empty 'to'")
	}
	if _, err := svc.SendText(context.Background(), "+5511999999999", ""); err == nil {
		t.Fatalf("expected validation error for empty 'body'")
	}
}

func TestMessagesService_SendText_GraphError(t *testing.T) {
	// Graph-style error payload
	errPayload := map[string]any{
		"error": map[string]any{
			"message":    "Invalid parameter",
			"type":       "OAuthException",
			"code":       100,
			"fbtrace_id": "abc123",
		},
	}
	b, _ := json.Marshal(errPayload)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(b)
	}))
	defer ts.Close()

	c := newTestClient(t, ts.URL)
	svc := services.NewMessagesService(c)
	_, err := svc.SendText(context.Background(), "+5511999999999", "hello")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	var ge *errorsx.GraphError
	if !errors.As(err, &ge) {
		t.Fatalf("expected GraphError, got %T: %v", err, err)
	}
	if ge.Detail.Code != 100 {
		t.Fatalf("unexpected graph code: %+v", ge.Detail)
	}
}

func TestMessagesService_SendText_Retry429Then200(t *testing.T) {
	fixturePath := filepath.Join("..", "..", "..", "testdata", "send_text_success.json")
	successBody, err := os.ReadFile(fixturePath)
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}

	var calls int
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		if calls == 1 {
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte(`{"error":{"message":"rate limit"}}`))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		io.Copy(w, bytes.NewReader(successBody))
	}))
	defer ts.Close()

	c := newTestClient(t, ts.URL)
	svc := services.NewMessagesService(c)
	resp, err := svc.SendText(context.Background(), "+5511999999999", "hello")
	if err != nil {
		t.Fatalf("SendText failed: %v", err)
	}
	if resp == nil || len(resp.Messages) == 0 {
		t.Fatalf("expected success response after retry, got %+v", resp)
	}
	if calls < 2 {
		t.Fatalf("expected at least 2 calls (429 then 200), got %d", calls)
	}
}
