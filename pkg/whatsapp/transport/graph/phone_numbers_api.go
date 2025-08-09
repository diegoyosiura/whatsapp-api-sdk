package graph

import (
	"context"
	"net/http"
)

// NewPhoneNumbersListRequest builds GET /{Version}/{WABA-ID}/phone_numbers
func NewPhoneNumbersListRequest(ctx context.Context, base, version, wabaID string) (*http.Request, error) {
	url := PhoneNumbersListEndpoint(base, version, wabaID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	return req, nil
}

// NewPhoneNumberGetRequest builds GET /{Version}/{Phone-Number-ID}
func NewPhoneNumberGetRequest(ctx context.Context, base, version, phoneNumberID string) (*http.Request, error) {
	url := PhoneNumberGetEndpoint(base, version, phoneNumberID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	return req, nil
}
