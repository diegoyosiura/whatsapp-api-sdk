package graph

import (
	"context"
	"fmt"
	"net/http"

	"github.com/diegoyosiura/whatsapp-sdk-go/pkg/whatsapp/domain"
	"github.com/diegoyosiura/whatsapp-sdk-go/pkg/whatsapp/ports"
)

type MediaAPI struct {
}

func (m *MediaAPI) GetMediaURL(ctx context.Context, base, version, mediaID string, tokenProvider ports.TokenProvider) (*http.Request, error) {
	url := RequestMediaLink(base, version, mediaID)

	token, err := tokenProvider.Token(ctx)
	if err != nil {
		return nil, fmt.Errorf("token: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	return req, nil
}

func (m *MediaAPI) Download(ctx context.Context, url *domain.DownloadLinkURL, tokenProvider ports.TokenProvider) (*http.Request, error) {
	token, err := tokenProvider.Token(ctx)
	if err != nil {
		return nil, fmt.Errorf("token: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url.Url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	return req, nil
}

func NewMediaAPI() ports.MediaAPI {
	return &MediaAPI{}
}
