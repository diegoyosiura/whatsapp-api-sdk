package domain

import "testing"

func TestParseWebhookEvent(t *testing.T) {
	ok := []byte(`{"object":"whatsapp","entry":[]}`)
	if _, err := ParseWebhookEvent(ok); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	bad := []byte("not-json")
	if _, err := ParseWebhookEvent(bad); err == nil {
		t.Fatalf("expected error for invalid json")
	}
}
