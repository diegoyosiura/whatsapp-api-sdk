package ports

import (
	"context"
	"net/http"

	"github.com/diegoyosiura/whatsapp-sdk-go/pkg/whatsapp/domain"
)

type MediaAPI interface {
	GetMediaURL(ctx context.Context, base, version, mediaID string, tokenProvider TokenProvider) (*http.Request, error)
	Download(ctx context.Context, url *domain.DownloadLinkURL, tokenProvider TokenProvider) (*http.Request, error)
}
