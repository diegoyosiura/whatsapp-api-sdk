package graph

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	portstesting "github.com/diegoyosiura/whatsapp-sdk-go/internal/testutils/whatsapp/ports"
)

func TestNewPhoneAPI_DefaultBase(t *testing.T) {
	a := NewPhoneAPI(&portstesting.FakeHTTPDoer{}, &portstesting.FakeTokenProvider{}, "v1", "waba", "")
	if a.baseURL != DefaultBaseURL {
		t.Fatalf("expected default base")
	}
}

func TestPhoneAPI_List_AttachAuthError(t *testing.T) {
	tp := &portstesting.FakeTokenProvider{Err: errors.New("no token")}
	a := NewPhoneAPI(&portstesting.FakeHTTPDoer{}, tp, "v1", "waba", "")
	if _, err := a.List(context.Background()); err == nil || !strings.Contains(err.Error(), "token:") {
		t.Fatalf("expected token error, got %v", err)
	}
}

func TestPhoneAPI_List_DoerError(t *testing.T) {
	doer := &portstesting.FakeHTTPDoer{Fn: func(ctx context.Context, req *http.Request) (*http.Response, error) {
		return nil, fmt.Errorf("boom")
	}}
	tp := &portstesting.FakeTokenProvider{TokenValue: "t"}
	a := NewPhoneAPI(doer, tp, "v1", "waba", "")
	if _, err := a.List(context.Background()); err == nil || !strings.Contains(err.Error(), "boom") {
		t.Fatalf("expected boom, got %v", err)
	}
}

func TestPhoneAPI_List_Non2xx(t *testing.T) {
	doer := &portstesting.FakeHTTPDoer{Fn: func(ctx context.Context, req *http.Request) (*http.Response, error) {
		resp := &http.Response{StatusCode: 500, Body: io.NopCloser(strings.NewReader(`{"error":{"message":"x"}}`)), Header: make(http.Header), Request: req}
		return resp, nil
	}}
	tp := &portstesting.FakeTokenProvider{TokenValue: "t"}
	a := NewPhoneAPI(doer, tp, "v1", "waba", "")
	if _, err := a.List(context.Background()); err == nil || !strings.Contains(err.Error(), "graph error") {
		t.Fatalf("expected graph error, got %v", err)
	}
}

func TestPhoneAPI_List_DecodeError(t *testing.T) {
	doer := &portstesting.FakeHTTPDoer{Fn: func(ctx context.Context, req *http.Request) (*http.Response, error) {
		resp := &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(`notjson`)), Header: make(http.Header), Request: req}
		return resp, nil
	}}
	tp := &portstesting.FakeTokenProvider{TokenValue: "t"}
	a := NewPhoneAPI(doer, tp, "v1", "waba", "")
	if _, err := a.List(context.Background()); err == nil || !strings.Contains(err.Error(), "decode phone list") {
		t.Fatalf("expected decode error, got %v", err)
	}
}

func TestPhoneAPI_List_Success(t *testing.T) {
	body := `{"data":[{"id":"1","display_phone_number":"1"}]}`
	doer := &portstesting.FakeHTTPDoer{Fn: func(ctx context.Context, req *http.Request) (*http.Response, error) {
		resp := &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(body)), Header: make(http.Header), Request: req}
		return resp, nil
	}}
	tp := &portstesting.FakeTokenProvider{TokenValue: "t"}
	a := NewPhoneAPI(doer, tp, "v1", "waba", "")
	out, err := a.List(context.Background())
	if err != nil || len(out.Data) != 1 || out.Data[0].ID != "1" {
		t.Fatalf("unexpected output %+v err %v", out, err)
	}
}

func TestPhoneAPI_Get_AttachAuthError(t *testing.T) {
	tp := &portstesting.FakeTokenProvider{Err: errors.New("no token")}
	a := NewPhoneAPI(&portstesting.FakeHTTPDoer{}, tp, "v1", "waba", "")
	if _, err := a.Get(context.Background(), "pn"); err == nil || !strings.Contains(err.Error(), "token:") {
		t.Fatalf("expected token error, got %v", err)
	}
}

func TestPhoneAPI_Get_DoerError(t *testing.T) {
	doer := &portstesting.FakeHTTPDoer{Fn: func(ctx context.Context, req *http.Request) (*http.Response, error) { return nil, fmt.Errorf("boom") }}
	tp := &portstesting.FakeTokenProvider{TokenValue: "t"}
	a := NewPhoneAPI(doer, tp, "v1", "waba", "")
	if _, err := a.Get(context.Background(), "pn"); err == nil || !strings.Contains(err.Error(), "boom") {
		t.Fatalf("expected boom, got %v", err)
	}
}

func TestPhoneAPI_Get_Non2xx(t *testing.T) {
	doer := &portstesting.FakeHTTPDoer{Fn: func(ctx context.Context, req *http.Request) (*http.Response, error) {
		resp := &http.Response{StatusCode: 404, Body: io.NopCloser(strings.NewReader(`{"error":{"message":"x"}}`)), Header: make(http.Header), Request: req}
		return resp, nil
	}}
	tp := &portstesting.FakeTokenProvider{TokenValue: "t"}
	a := NewPhoneAPI(doer, tp, "v1", "waba", "")
	if _, err := a.Get(context.Background(), "pn"); err == nil || !strings.Contains(err.Error(), "graph error") {
		t.Fatalf("expected graph error, got %v", err)
	}
}

func TestPhoneAPI_Get_DecodeError(t *testing.T) {
	doer := &portstesting.FakeHTTPDoer{Fn: func(ctx context.Context, req *http.Request) (*http.Response, error) {
		resp := &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(`bad`)), Header: make(http.Header), Request: req}
		return resp, nil
	}}
	tp := &portstesting.FakeTokenProvider{TokenValue: "t"}
	a := NewPhoneAPI(doer, tp, "v1", "waba", "")
	if _, err := a.Get(context.Background(), "pn"); err == nil || !strings.Contains(err.Error(), "decode phone") {
		t.Fatalf("expected decode error, got %v", err)
	}
}

func TestPhoneAPI_Get_Success(t *testing.T) {
	body := `{"id":"pn","display_phone_number":"1"}`
	doer := &portstesting.FakeHTTPDoer{Fn: func(ctx context.Context, req *http.Request) (*http.Response, error) {
		resp := &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(body)), Header: make(http.Header), Request: req}
		return resp, nil
	}}
	tp := &portstesting.FakeTokenProvider{TokenValue: "t"}
	a := NewPhoneAPI(doer, tp, "v1", "waba", "")
	out, err := a.Get(context.Background(), "pn")
	if err != nil || out.ID != "pn" {
		t.Fatalf("unexpected output %v err %v", out, err)
	}
}
