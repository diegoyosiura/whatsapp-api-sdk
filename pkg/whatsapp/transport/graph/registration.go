package graph

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/diegoyosiura/whatsapp-sdk-go/pkg/errorsx"
	"github.com/diegoyosiura/whatsapp-sdk-go/pkg/whatsapp/domain"
	"github.com/diegoyosiura/whatsapp-sdk-go/pkg/whatsapp/ports"
)

// RegistrationAPI is a thin Graph adapter that maps domain requests to Graph API endpoints.
// It strictly depends on hex ports and never on service layer types.
type RegistrationAPI struct {
	doer          ports.HTTPDoer
	tokenProvider ports.TokenProvider
	version       string
	phoneNumberID string
	baseURL       string // default: https://graph.facebook.com
}

// compile-time check
var _ ports.RegistrationAPI = (*RegistrationAPI)(nil)

// NewRegistrationAPI wires the adapter with hexagonal dependencies and static params.
func NewRegistrationAPI(doer ports.HTTPDoer, token ports.TokenProvider, version, phoneNumberID, baseURL string) *RegistrationAPI {
	if baseURL == "" {
		baseURL = DefaultBaseURL
	}
	return &RegistrationAPI{
		doer:          doer,
		tokenProvider: token,
		version:       version,
		phoneNumberID: phoneNumberID,
		baseURL:       baseURL,
	}
}

func (a *RegistrationAPI) endpoint(path string) string {
	// {base}/{version}/{phoneNumberID}/{path}
	if path == "" {
		return fmt.Sprintf("%s/%s/%s", a.baseURL, a.version, a.phoneNumberID)
	}
	return fmt.Sprintf("%s/%s/%s/%s", a.baseURL, a.version, a.phoneNumberID, path)
}

// RequestCode -> POST /{Version}/{Phone-Number-ID}/request_code
func (a *RegistrationAPI) RequestCode(ctx context.Context, p domain.RequestCodeParams) (*domain.ActionResult, error) {
	body, _ := json.Marshal(p)
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, a.endpoint("request_code"), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	if err := a.attachAuth(ctx, req); err != nil {
		return nil, err
	}
	return a.decodeActionResult(req)
}

// VerifyCode -> POST /{Version}/{Phone-Number-ID}/verify_code
func (a *RegistrationAPI) VerifyCode(ctx context.Context, p domain.VerifyCodeParams) (*domain.ActionResult, error) {
	body, _ := json.Marshal(p)
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, a.endpoint("verify_code"), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	if err := a.attachAuth(ctx, req); err != nil {
		return nil, err
	}
	return a.decodeActionResult(req)
}

// Register -> POST /{Version}/{Phone-Number-ID}/register
// Graph requires { "messaging_product":"whatsapp", "pin":"XXXXXX"? }
func (a *RegistrationAPI) Register(ctx context.Context, p domain.RegisterParams) (*domain.ActionResult, error) {
	payload := map[string]any{
		"messaging_product": "whatsapp",
	}
	if p.Pin != nil && *p.Pin != "" {
		payload["pin"] = *p.Pin
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, a.endpoint("register"), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	if err := a.attachAuth(ctx, req); err != nil {
		return nil, err
	}
	return a.decodeActionResult(req)
}

// Deregister -> POST /{Version}/{Phone-Number-ID}/deregister
func (a *RegistrationAPI) Deregister(ctx context.Context) (*domain.ActionResult, error) {
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, a.endpoint("deregister"), http.NoBody)
	req.Header.Set("Content-Type", "application/json")
	if err := a.attachAuth(ctx, req); err != nil {
		return nil, err
	}
	return a.decodeActionResult(req)
}

// SetTwoStep -> POST /{Version}/{Phone-Number-ID}  with { "pin": "XXXXXX" }
func (a *RegistrationAPI) SetTwoStep(ctx context.Context, p domain.TwoStepParams) (*domain.ActionResult, error) {
	body, _ := json.Marshal(map[string]string{"pin": p.Pin})
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, a.endpoint(""), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	if err := a.attachAuth(ctx, req); err != nil {
		return nil, err
	}
	return a.decodeActionResult(req)
}

// attachAuth injects the Bearer access token from TokenProvider.
// TODO: align method name/signature with your ports.TokenProvider.
func (a *RegistrationAPI) attachAuth(ctx context.Context, req *http.Request) error {
	// Example expectation:
	//   token, err := a.tokenProvider.Token(ctx)
	// Adjust if your interface is Get(ctx) or Provide(ctx).
	token, err := a.tokenProvider.Token(ctx)
	if err != nil {
		return fmt.Errorf("token: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	return nil
}

func (a *RegistrationAPI) decodeActionResult(req *http.Request) (*domain.ActionResult, error) {
	resp, err := a.doer.Do(req.Context(), req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// The Graph API normally replies { "success": true } for these endpoints.
	dec := json.NewDecoder(resp.Body)
	var out domain.ActionResult
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		if err := dec.Decode(&out); err != nil {
			return nil, fmt.Errorf("decode success: %w", err)
		}
		return &out, nil
	}

	// Non-2xx: try to decode Graph error envelope, or fall back.
	// TODO: if you have a GraphError type + DecodeGraphError, use it here.
	body, _ := io.ReadAll(resp.Body)
	if ge := errorsx.TryParseGraphError(resp, body); ge != nil {
		return nil, ge
	}
	return nil, errorsx.NewHTTPErrorFromResponse(resp, body)
}
