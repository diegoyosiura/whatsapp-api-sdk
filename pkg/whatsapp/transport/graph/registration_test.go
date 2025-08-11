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
	"github.com/diegoyosiura/whatsapp-sdk-go/pkg/whatsapp/domain"
)

func TestNewRegistrationAPI_DefaultBase(t *testing.T) {
	a := NewRegistrationAPI(&portstesting.FakeHTTPDoer{}, &portstesting.FakeTokenProvider{}, "v1", "pn", "")
	if a.baseURL != DefaultBaseURL {
		t.Fatalf("expected default base")
	}
}

func TestRegistrationAPI_endpoint(t *testing.T) {
	a := NewRegistrationAPI(nil, nil, "v1", "pn", "https://g")
	if got := a.endpoint("register"); got != "https://g/v1/pn/register" {
		t.Fatalf("endpoint wrong: %s", got)
	}
	if got := a.endpoint(""); got != "https://g/v1/pn" {
		t.Fatalf("endpoint empty wrong: %s", got)
	}
}

func TestRegistrationAPI_RequestCode_AttachAuthErr(t *testing.T) {
	tp := &portstesting.FakeTokenProvider{Err: errors.New("no token")}
	a := NewRegistrationAPI(&portstesting.FakeHTTPDoer{}, tp, "v1", "pn", "")
	_, err := a.RequestCode(context.Background(), domain.RequestCodeParams{CodeMethod: domain.CodeMethodSMS})
	if err == nil || !strings.Contains(err.Error(), "token:") {
		t.Fatalf("expected token error, got %v", err)
	}
}

func TestRegistrationAPI_VerifyCode_AttachAuthErr(t *testing.T) {
	tp := &portstesting.FakeTokenProvider{Err: errors.New("no token")}
	a := NewRegistrationAPI(&portstesting.FakeHTTPDoer{}, tp, "v1", "pn", "")
	_, err := a.VerifyCode(context.Background(), domain.VerifyCodeParams{Code: "1"})
	if err == nil || !strings.Contains(err.Error(), "token:") {
		t.Fatalf("expected token error, got %v", err)
	}
}

func TestRegistrationAPI_Register_AttachAuthErr(t *testing.T) {
	tp := &portstesting.FakeTokenProvider{Err: errors.New("no token")}
	a := NewRegistrationAPI(&portstesting.FakeHTTPDoer{}, tp, "v1", "pn", "")
	_, err := a.Register(context.Background(), domain.RegisterParams{})
	if err == nil || !strings.Contains(err.Error(), "token:") {
		t.Fatalf("expected token error, got %v", err)
	}
}

func TestRegistrationAPI_Deregister_AttachAuthErr(t *testing.T) {
	tp := &portstesting.FakeTokenProvider{Err: errors.New("no token")}
	a := NewRegistrationAPI(&portstesting.FakeHTTPDoer{}, tp, "v1", "pn", "")
	_, err := a.Deregister(context.Background())
	if err == nil || !strings.Contains(err.Error(), "token:") {
		t.Fatalf("expected token error, got %v", err)
	}
}

func TestRegistrationAPI_SetTwoStep_AttachAuthErr(t *testing.T) {
	tp := &portstesting.FakeTokenProvider{Err: errors.New("no token")}
	a := NewRegistrationAPI(&portstesting.FakeHTTPDoer{}, tp, "v1", "pn", "")
	_, err := a.SetTwoStep(context.Background(), domain.TwoStepParams{Pin: "1"})
	if err == nil || !strings.Contains(err.Error(), "token:") {
		t.Fatalf("expected token error, got %v", err)
	}
}

func newRegistrationAPISuccess(fn func(ctx context.Context, req *http.Request) (*http.Response, error)) *RegistrationAPI {
	doer := &portstesting.FakeHTTPDoer{Fn: fn}
	tp := &portstesting.FakeTokenProvider{TokenValue: "t"}
	return NewRegistrationAPI(doer, tp, "v1", "pn", "")
}

func successResponse(req *http.Request, body string) *http.Response {
	return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(body)), Header: make(http.Header), Request: req}
}

func TestRegistrationAPI_RequestCode_Success(t *testing.T) {
	a := newRegistrationAPISuccess(func(ctx context.Context, req *http.Request) (*http.Response, error) {
		return successResponse(req, `{"success":true}`), nil
	})
	out, err := a.RequestCode(context.Background(), domain.RequestCodeParams{CodeMethod: domain.CodeMethodSMS})
	if err != nil || !out.Success {
		t.Fatalf("unexpected %v %v", out, err)
	}
}

func TestRegistrationAPI_VerifyCode_Success(t *testing.T) {
	a := newRegistrationAPISuccess(func(ctx context.Context, req *http.Request) (*http.Response, error) {
		return successResponse(req, `{"success":true}`), nil
	})
	out, err := a.VerifyCode(context.Background(), domain.VerifyCodeParams{Code: "123"})
	if err != nil || !out.Success {
		t.Fatalf("unexpected %v %v", out, err)
	}
}

func TestRegistrationAPI_Register_Success(t *testing.T) {
	a := newRegistrationAPISuccess(func(ctx context.Context, req *http.Request) (*http.Response, error) {
		b, _ := io.ReadAll(req.Body)
		if !strings.Contains(string(b), "pin") {
			t.Fatalf("expected pin in body: %s", string(b))
		}
		return successResponse(req, `{"success":true}`), nil
	})
	pin := "1234"
	out, err := a.Register(context.Background(), domain.RegisterParams{Pin: &pin})
	if err != nil || !out.Success {
		t.Fatalf("unexpected %v %v", out, err)
	}
}

func TestRegistrationAPI_Deregister_Success(t *testing.T) {
	a := newRegistrationAPISuccess(func(ctx context.Context, req *http.Request) (*http.Response, error) {
		return successResponse(req, `{"success":true}`), nil
	})
	out, err := a.Deregister(context.Background())
	if err != nil || !out.Success {
		t.Fatalf("unexpected %v %v", out, err)
	}
}

func TestRegistrationAPI_SetTwoStep_Success(t *testing.T) {
	a := newRegistrationAPISuccess(func(ctx context.Context, req *http.Request) (*http.Response, error) {
		return successResponse(req, `{"success":true}`), nil
	})
	out, err := a.SetTwoStep(context.Background(), domain.TwoStepParams{Pin: "222"})
	if err != nil || !out.Success {
		t.Fatalf("unexpected %v %v", out, err)
	}
}

func TestRegistrationAPI_decodeActionResult_DoerError(t *testing.T) {
	doer := &portstesting.FakeHTTPDoer{Fn: func(ctx context.Context, req *http.Request) (*http.Response, error) {
		return nil, fmt.Errorf("boom")
	}}
	a := NewRegistrationAPI(doer, &portstesting.FakeTokenProvider{}, "v1", "pn", "")
	req, _ := http.NewRequestWithContext(context.Background(), http.MethodPost, "https://x", nil)
	if _, err := a.decodeActionResult(req); err == nil || !strings.Contains(err.Error(), "boom") {
		t.Fatalf("expected boom, got %v", err)
	}
}

func TestRegistrationAPI_decodeActionResult_DecodeError(t *testing.T) {
	doer := &portstesting.FakeHTTPDoer{Fn: func(ctx context.Context, req *http.Request) (*http.Response, error) {
		return successResponse(req, `bad`), nil
	}}
	a := NewRegistrationAPI(doer, &portstesting.FakeTokenProvider{}, "v1", "pn", "")
	req, _ := http.NewRequestWithContext(context.Background(), http.MethodPost, "https://x", nil)
	if _, err := a.decodeActionResult(req); err == nil || !strings.Contains(err.Error(), "decode success") {
		t.Fatalf("expected decode error, got %v", err)
	}
}

func TestRegistrationAPI_decodeActionResult_GraphError(t *testing.T) {
	doer := &portstesting.FakeHTTPDoer{Fn: func(ctx context.Context, req *http.Request) (*http.Response, error) {
		resp := &http.Response{StatusCode: 400, Body: io.NopCloser(strings.NewReader(`{"error":{"message":"x"}}`)), Header: make(http.Header), Request: req}
		return resp, nil
	}}
	a := NewRegistrationAPI(doer, &portstesting.FakeTokenProvider{}, "v1", "pn", "")
	req, _ := http.NewRequestWithContext(context.Background(), http.MethodPost, "https://x", nil)
	if _, err := a.decodeActionResult(req); err == nil || !strings.Contains(err.Error(), "graph error") {
		t.Fatalf("expected graph error, got %v", err)
	}
}
