package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/diegoyosiura/whatsapp-sdk-go/pkg/whatsapp/domain"
	"github.com/diegoyosiura/whatsapp-sdk-go/pkg/whatsapp/ports"
	"github.com/diegoyosiura/whatsapp-sdk-go/pkg/whatsapp/transport/graph"
)

type MediaService struct {
	api ports.MediaAPI
	c   clientCore
}

func NewMediaService(c clientCore) *MediaService {
	return &MediaService{api: graph.NewMediaAPI(), c: c}
}

func (s *MediaService) UploadMedia(ctx context.Context, f *os.File) (*domain.MediaUpload, error) {
	rq, err := s.api.Upload(ctx, s.c.BaseURL(), s.c.Version(), s.c.PhoneNumberID(), f, s.c.TokenProvider())
	if err != nil {
		return nil, err
	}

	resp, err := s.c.Do(ctx, rq)

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 8<<10))
		return nil, fmt.Errorf("bad status: %s - %s", resp.Status, string(b))
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return nil, fmt.Errorf("read canceled: %w", err)
		}
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, fmt.Errorf("read timeout: %w", err)
		}
		return nil, fmt.Errorf("read response body: %w", err)
	}

	response := &domain.MediaUpload{}
	if err := json.Unmarshal(b, response); err != nil {
		return nil, err
	}
	return response, nil
}

func (s *MediaService) DeleteMedia(ctx context.Context, mediaID string) (*domain.MediaDelete, error) {
	rq, err := s.api.Delete(ctx, s.c.BaseURL(), s.c.Version(), s.c.PhoneNumberID(), mediaID, s.c.TokenProvider())
	if err != nil {
		return nil, err
	}

	resp, err := s.c.Do(ctx, rq)

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 8<<10))
		return nil, fmt.Errorf("bad status: %s - %s", resp.Status, string(b))
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return nil, fmt.Errorf("read canceled: %w", err)
		}
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, fmt.Errorf("read timeout: %w", err)
		}
		return nil, fmt.Errorf("read response body: %w", err)
	}

	response := &domain.MediaDelete{}
	if err := json.Unmarshal(b, response); err != nil {
		return nil, err
	}
	return response, nil
}
func (s *MediaService) GetMediaURL(ctx context.Context, mediaID string) (*domain.DownloadLinkURL, error) {
	rq, err := s.api.GetMediaURL(ctx, s.c.BaseURL(), s.c.Version(), mediaID, s.c.TokenProvider())
	if err != nil {
		return nil, err
	}

	resp, err := s.c.Do(ctx, rq) // use gctx here
	if err != nil {
		return nil, err
	}
	defer func() {
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 8<<10))
		return nil, fmt.Errorf("get media url bad status: %s - %s", resp.Status, string(b))
	}

	d := &domain.DownloadLinkURL{}
	if err := json.NewDecoder(resp.Body).Decode(d); err != nil {
		return nil, fmt.Errorf("decode request: %w", err)
	}

	return d, nil
}

func (s *MediaService) DownloadMedia(ctx context.Context, d *domain.DownloadLinkURL, fm ports.FileManagerAPI) (ports.FileManagerAPI, error) {
	base := context.WithoutCancel(ctx)
	dctx, cancel := context.WithTimeout(base, 30*time.Second)
	defer cancel()

	rq, err := s.api.Download(dctx, d, s.c.TokenProvider())
	if err != nil {
		return nil, err
	}

	resp, err := s.c.Do(dctx, rq) // use dctx aqui
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return nil, fmt.Errorf("download canceled: %w", err)
		}
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, fmt.Errorf("download timeout: %w", err)
		}
		return nil, err
	}
	defer func() {
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 8<<10))
		return nil, fmt.Errorf("bad status: %s - %s", resp.Status, string(b))
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return nil, fmt.Errorf("read canceled: %w", err)
		}
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, fmt.Errorf("read timeout: %w", err)
		}
		return nil, fmt.Errorf("read response body: %w", err)
	}

	if err := fm.SetData(dctx, b); err != nil { // use dctx aqui também
		return nil, fmt.Errorf("set data: %w", err)
	}
	return fm, nil
}
func (s *MediaService) GetMedia(ctx context.Context, message domain.InboundMessage, fm ports.FileManagerAPI) (ports.FileManagerAPI, error) {
	switch message.Type {
	case "audio":
		// Detached context with a short timeout for metadata call
		gbase := context.WithoutCancel(ctx)
		gctx, gcancel := context.WithTimeout(gbase, 10*time.Second)
		defer gcancel()

		d, err := s.GetMediaURL(gctx, message.Audio.ID)
		if err != nil {
			return nil, err
		}

		return s.DownloadMedia(ctx, d, fm) // DownloadMedia criará seu próprio dctx
	default:
		return nil, errors.New("not supported")
	}
}
