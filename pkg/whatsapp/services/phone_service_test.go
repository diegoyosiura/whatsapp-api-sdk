package services_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/diegoyosiura/whatsapp-sdk-go/pkg/errorsx"
	"github.com/diegoyosiura/whatsapp-sdk-go/pkg/whatsapp/domain"
	"github.com/diegoyosiura/whatsapp-sdk-go/pkg/whatsapp/ports"
	"github.com/diegoyosiura/whatsapp-sdk-go/pkg/whatsapp/services"
)

// fakePhoneAPI implements ports.PhoneAPI for unit testing the service.
type fakePhoneAPI struct {
	listFn func() (*domain.PhoneList, error)
	getFn  func(id string) (*domain.Phone, error)
}

func (f *fakePhoneAPI) List(_ context.Context) (*domain.PhoneList, error) { return f.listFn() }
func (f *fakePhoneAPI) Get(_ context.Context, id string) (*domain.Phone, error) {
	return f.getFn(id)
}

func TestPhoneService_List_Success(t *testing.T) {
	fixturePath := filepath.Join("..", "..", "..", "testdata", "phone_numbers_list.json")
	b, err := os.ReadFile(fixturePath)
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}

	var out domain.PhoneList
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	var api ports.PhoneAPI = &fakePhoneAPI{
		listFn: func() (*domain.PhoneList, error) { return &out, nil },
		getFn:  nil,
	}
	svc := services.NewPhoneService(api)
	resp, err := svc.List(context.Background())

	if err != nil {
		t.Fatalf("List: unexpected error: %v", err)
	}
	if resp == nil || len(resp.Data) == 0 {
		t.Fatalf("expected at least one phone number, got %v", resp)
	}
}

func TestPhoneService_Get_Success(t *testing.T) {
	fixturePath := filepath.Join("..", "..", "..", "testdata", "phone_number_get.json")
	b, err := os.ReadFile(fixturePath)
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}
	var p domain.Phone
	if err := json.Unmarshal(b, &p); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	api := &fakePhoneAPI{
		listFn: nil,
		getFn:  func(id string) (*domain.Phone, error) { return &p, nil },
	}
	svc := services.NewPhoneService(api)
	resp, err := svc.Get(context.Background(), "1234567890")
	if err != nil {
		t.Fatalf("Get: unexpected error: %v", err)
	}
	if resp == nil || resp.ID == "" {
		t.Fatalf("expected phone number details, got %+v", resp)
	}
}

func TestPhoneService_List_GraphError(t *testing.T) {
	errPayload := map[string]any{
		"error": map[string]any{
			"message":    "Invalid parameter",
			"type":       "OAuthException",
			"code":       100,
			"fbtrace_id": "xyz",
		},
	}
	b, _ := json.Marshal(errPayload)

	_ = b
	api := &fakePhoneAPI{
		listFn: func() (*domain.PhoneList, error) {
			// Monte uma *http.Response "real" para o builder de erro não dereferenciar nil.
			u, _ := url.Parse("https://graph.facebook.com/v20.0/123/phone_numbers")
			resp := &http.Response{
				StatusCode: http.StatusBadRequest,
				Status:     http.StatusText(http.StatusBadRequest),
				Header:     make(http.Header),
				Body:       io.NopCloser(bytes.NewReader(b)),
				Request:    &http.Request{Method: http.MethodGet, URL: u},
			}
			resp.Header.Set("Content-Type", "application/json")

			he := errorsx.NewHTTPErrorFromResponse(resp, b)
			return nil, &errorsx.GraphError{HTTP: he, Raw: b}
		},
	}
	svc := services.NewPhoneService(api)

	_, err := svc.List(context.Background())
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	var ge *errorsx.GraphError
	if !errors.As(err, &ge) {
		t.Fatalf("expected GraphError, got %T", err)
	}
}
