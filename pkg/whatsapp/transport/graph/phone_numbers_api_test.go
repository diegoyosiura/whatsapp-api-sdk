package graph

import (
	"context"
	"net/http"
	"testing"
)

func TestPhoneNumbersRequests(t *testing.T) {
	ctx := context.Background()
	req, err := NewPhoneNumbersListRequest(ctx, DefaultBaseURL, "v1", "waba")
	if err != nil {
		t.Fatalf("list request err: %v", err)
	}
	if req.Method != http.MethodGet || req.URL.String() != PhoneNumbersListEndpoint(DefaultBaseURL, "v1", "waba") {
		t.Fatalf("unexpected request: %s %s", req.Method, req.URL.String())
	}
	if req.Header.Get("Accept") != "application/json" {
		t.Fatalf("missing accept header")
	}

	req2, err := NewPhoneNumberGetRequest(ctx, DefaultBaseURL, "v1", "pn")
	if err != nil {
		t.Fatalf("get request err: %v", err)
	}
	if req2.Method != http.MethodGet || req2.URL.String() != PhoneNumberGetEndpoint(DefaultBaseURL, "v1", "pn") {
		t.Fatalf("unexpected request2: %s %s", req2.Method, req2.URL.String())
	}
	if req2.Header.Get("Accept") != "application/json" {
		t.Fatalf("missing accept header on req2")
	}
}
