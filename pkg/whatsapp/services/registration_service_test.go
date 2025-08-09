package services_test

import (
	"context"
	"encoding/json"
	"errors"

	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/diegoyosiura/whatsapp-sdk-go/pkg/errorsx"
	"github.com/diegoyosiura/whatsapp-sdk-go/pkg/whatsapp"
	"github.com/diegoyosiura/whatsapp-sdk-go/pkg/whatsapp/domain"
	"github.com/diegoyosiura/whatsapp-sdk-go/pkg/whatsapp/services"
)

type fakeRegAPI struct{}

func (a fakeRegAPI) Deregister(ctx context.Context) (*domain.ActionResult, error) {
	//TODO implement me
	panic("implement me")
}

func (a fakeRegAPI) SetTwoStep(ctx context.Context, p domain.TwoStepParams) (*domain.ActionResult, error) {
	//TODO implement me
	panic("implement me")
}

func (fakeRegAPI) RequestCode(ctx context.Context, p domain.RequestCodeParams) (*domain.ActionResult, error) {
	return &domain.ActionResult{Success: true}, nil
}
func (fakeRegAPI) VerifyCode(ctx context.Context, p domain.VerifyCodeParams) (*domain.ActionResult, error) {
	return &domain.ActionResult{Success: true}, nil
}
func (fakeRegAPI) Register(ctx context.Context, p domain.RegisterParams) (*domain.ActionResult, error) {
	return &domain.ActionResult{Success: true}, nil
}

type regFakeToken struct{}

func (regFakeToken) Token(ctx context.Context) (string, error) { return "t", nil }
func (regFakeToken) Refresh(ctx context.Context) error         { return nil }

func newRegClient(t *testing.T, baseURL string) *whatsapp.Client {
	t.Helper()
	c, err := whatsapp.NewClient(whatsapp.Options{
		Version:       "v20.0",
		WABAID:        "waba",
		PhoneNumberID: "pnid",
		TokenProvider: regFakeToken{},
		BaseURL:       baseURL,
		Timeout:       1 * time.Second,
		RetryMax:      0,
		UserAgent:     "test",
	})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	return c
}

func TestRegistration_RequestVerifyRegister_OK(t *testing.T) {
	svc := services.NewRegistrationService(fakeRegAPI{})

	if out, err := svc.RequestCode(context.Background(), domain.RequestCodeParams{CodeMethod: domain.CodeMethodSMS, Locale: "en_US"}); err != nil || !out.Success {
		t.Fatalf("RequestCode err=%v out=%+v", err, out)
	}
	if out, err := svc.VerifyCode(context.Background(), domain.VerifyCodeParams{Code: "123456"}); err != nil || !out.Success {
		t.Fatalf("VerifyCode err=%v out=%+v", err, out)
	}
	pin := "123456"
	if out, err := svc.Register(context.Background(), domain.RegisterParams{Pin: &pin}); err != nil || !out.Success {
		t.Fatalf("Register err=%v out=%+v", err, out)
	}
}

func TestRegistration_ErrorGraphPayload(t *testing.T) {
	// Responde 400 com corpo de erro Graph para /register
	mux := http.NewServeMux()
	mux.HandleFunc("/v20.0/pnid/register", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		errBody := map[string]any{
			"error": map[string]any{
				"message": "Invalid PIN",
				"type":    "OAuthException",
				"code":    190,
			},
		}
		b, _ := json.Marshal(errBody)
		w.Write(b)
	})
	ts := httptest.NewServer(mux)
	defer ts.Close()

	c := newRegClient(t, ts.URL)
	svc := services.NewRegistrationService(c.RegistrationService().API())

	_, err := svc.Register(context.Background(), domain.RegisterParams{Pin: ptr("000000")})
	if err == nil {
		t.Fatalf("expected error")
	}
	var ge *errorsx.GraphError
	if ok := As(err, &ge); !ok {
		t.Fatalf("expected GraphError, got %T", err)
	}
}

func ptr(s string) *string { return &s }

// small helper: errors.As across go versions used in tests only
func As(err error, target any) bool { return errorsxAs(err, target) }

// shim to avoid importing "errors" in multiple places
func errorsxAs(err error, target any) bool {
	type as interface{ As(error, any) bool }
	// Use std errors.As via a tiny indirection to keep this file focused
	return func(e error, t any) bool {
		return (func() bool {
			// inline: stdlib
			return errors.As(e, &t)
		})()
	}(err, target)
}
