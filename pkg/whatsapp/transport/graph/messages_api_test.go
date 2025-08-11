package graph

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	portstesting "github.com/diegoyosiura/whatsapp-sdk-go/internal/testutils/whatsapp/ports"
)

func TestNewSendTextHTTPRequest(t *testing.T) {
	ctx := context.Background()
	payload := SendTextRequest{MessagingProduct: "whatsapp", To: "123", Type: "text"}
	payload.Text.Body = "hi"
	tp := &portstesting.FakeTokenProvider{TokenValue: "tok"}
	req, err := NewSendTextHTTPRequest(ctx, DefaultBaseURL, "v1", "pn", payload, tp)
	if err != nil {
		t.Fatalf("NewSendTextHTTPRequest error: %v", err)
	}
	if req.Method != http.MethodPost {
		t.Fatalf("expected POST got %s", req.Method)
	}
	if req.URL.String() != MessagesEndpoint(DefaultBaseURL, "v1", "pn") {
		t.Fatalf("unexpected url %s", req.URL.String())
	}
	if got := req.Header.Get("Authorization"); got != "Bearer tok" {
		t.Fatalf("auth header %q", got)
	}
	b, _ := io.ReadAll(req.Body)
	if !strings.Contains(string(b), "\"hi\"") {
		t.Fatalf("body not encoded: %s", string(b))
	}
	rc, err := req.GetBody()
	if err != nil {
		t.Fatalf("GetBody error: %v", err)
	}
	b2, _ := io.ReadAll(rc)
	if string(b2) != string(b) {
		t.Fatalf("GetBody mismatch")
	}
}

func TestNewSendTextHTTPRequestTokenError(t *testing.T) {
	ctx := context.Background()
	payload := SendTextRequest{MessagingProduct: "whatsapp", To: "123", Type: "text"}
	tp := &portstesting.FakeTokenProvider{Err: fmt.Errorf("no token")}
	_, err := NewSendTextHTTPRequest(ctx, DefaultBaseURL, "v1", "pn", payload, tp)
	if err == nil || !strings.Contains(err.Error(), "token:") {
		t.Fatalf("expected token error, got %v", err)
	}
}
