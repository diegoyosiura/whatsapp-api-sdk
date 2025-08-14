package ports

import (
	"context"
	"net/http"
	"os"

	"github.com/diegoyosiura/whatsapp-sdk-go/pkg/whatsapp/domain"
)

type MediaAPI interface {
	Upload(ctx context.Context, base, version, phoneNumberID string, f *os.File, tokenProvider TokenProvider) (*http.Request, error)
	Download(ctx context.Context, url *domain.DownloadLinkURL, tokenProvider TokenProvider) (*http.Request, error)
	Delete(ctx context.Context, base, version, phoneNumberID, mediaID string, tokenProvider TokenProvider) (*http.Request, error)
	GetMediaURL(ctx context.Context, base, version, mediaID string, tokenProvider TokenProvider) (*http.Request, error)
}
