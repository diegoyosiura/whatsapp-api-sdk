package graph

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/diegoyosiura/whatsapp-sdk-go/pkg/whatsapp/domain"
	"github.com/diegoyosiura/whatsapp-sdk-go/pkg/whatsapp/ports"
)

type PhoneAPI struct {
	doer          ports.HTTPDoer
	tokenProvider ports.TokenProvider
	version       string
	wabaID        string
	baseURL       string // default: https://graph.facebook.com
}

func NewPhoneAPI(doer ports.HTTPDoer, token ports.TokenProvider, version, wabaID string) *PhoneAPI {
	return &PhoneAPI{
		doer:          doer,
		tokenProvider: token,
		version:       version,
		wabaID:        wabaID,
		baseURL:       "https://graph.facebook.com",
	}
}

var _ ports.PhoneAPI = (*PhoneAPI)(nil)

func (a *PhoneAPI) endpointWABA(path string) string {
	return fmt.Sprintf("%s/%s/%s/%s", a.baseURL, a.version, a.wabaID, path)
}

func (a *PhoneAPI) endpointPhoneID(phoneID string) string {
	return fmt.Sprintf("%s/%s/%s", a.baseURL, a.version, phoneID)
}

func (a *PhoneAPI) attachAuth(ctx context.Context, req *http.Request) error {
	token, err := a.tokenProvider.Token(ctx)
	if err != nil {
		return fmt.Errorf("token: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	return nil
}

func (a *PhoneAPI) List(ctx context.Context) (*domain.PhoneList, error) {
	u := a.endpointWABA("phone_numbers")
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, u, http.NoBody)
	if err := a.attachAuth(ctx, req); err != nil {
		return nil, err
	}

	resp, err := a.doer.Do(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		var any map[string]any
		_ = json.NewDecoder(resp.Body).Decode(&any)
		return nil, fmt.Errorf("graph error %d: %v", resp.StatusCode, any)
	}

	var out domain.PhoneList
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("decode phone list: %w", err)
	}
	return &out, nil
}

func (a *PhoneAPI) Get(ctx context.Context, phoneID string) (*domain.Phone, error) {
	// Ask for the common fields to keep output consistent with List.
	q := url.Values{}
	q.Set("fields", "id,display_phone_number,verified_name,quality_rating,is_official_business_account,account_mode")

	u := a.endpointPhoneID(phoneID) + "?" + q.Encode()
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, u, http.NoBody)
	if err := a.attachAuth(ctx, req); err != nil {
		return nil, err
	}

	resp, err := a.doer.Do(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		var any map[string]any
		_ = json.NewDecoder(resp.Body).Decode(&any)
		return nil, fmt.Errorf("graph error %d: %v", resp.StatusCode, any)
	}

	var out domain.Phone
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("decode phone: %w", err)
	}
	return &out, nil
}
