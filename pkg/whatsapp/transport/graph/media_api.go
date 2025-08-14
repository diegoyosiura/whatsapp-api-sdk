package graph

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/diegoyosiura/whatsapp-sdk-go/pkg/utils"
	"github.com/diegoyosiura/whatsapp-sdk-go/pkg/whatsapp/domain"
	"github.com/diegoyosiura/whatsapp-sdk-go/pkg/whatsapp/ports"
)

type MediaAPI struct {
}

func NewMediaAPI() ports.MediaAPI {
	return &MediaAPI{}
}

func (m *MediaAPI) Upload(ctx context.Context, base, version, phoneNumberID string, f *os.File, tokenProvider ports.TokenProvider) (*http.Request, error) {
	defer func() {
		_ = f.Close()
	}()

	token, err := tokenProvider.Token(ctx)
	if err != nil {
		return nil, fmt.Errorf("token: %w", err)
	}

	url := RequestMediaUpload(base, version, phoneNumberID)
	buf := &bytes.Buffer{}
	writer := multipart.NewWriter(buf)

	_ = writer.WriteField("messaging_product", "whatsapp")

	part, err := writer.CreatePart(utils.GetMimeTypeHeader(f))
	if err != nil {
		_ = writer.Close()
		return nil, fmt.Errorf("writer: %w", err)
	}
	if _, err := io.Copy(part, f); err != nil {
		_ = writer.Close()
		return nil, err
	}

	_ = writer.Close()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", writer.FormDataContentType())

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

func (m *MediaAPI) Delete(ctx context.Context, base, version, phoneNumberID, mediaID string, tokenProvider ports.TokenProvider) (*http.Request, error) {
	token, err := tokenProvider.Token(ctx)
	if err != nil {
		return nil, fmt.Errorf("token: %w", err)
	}

	url := RequestMediaDelete(base, version, mediaID)
	url = fmt.Sprintf("%s?phone_number_id=%s", url, phoneNumberID)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	return req, nil
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
