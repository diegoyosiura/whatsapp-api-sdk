package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
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
func (s *MediaService) GetMedia(ctx context.Context, message domain.InboundMessage, fm ports.FileManagerAPI) (ports.FileManagerAPI, error) {
	switch message.Type {
	case "audio":
		// Detached context with a short timeout for metadata call
		gbase := context.WithoutCancel(ctx)
		gctx, gcancel := context.WithTimeout(gbase, 10*time.Second)
		defer gcancel()

		rq, err := s.api.GetMediaURL(gctx, s.c.BaseURL(), s.c.Version(), message.Audio.ID, s.c.TokenProvider())
		if err != nil {
			return nil, err
		}

		resp, err := s.c.Do(gctx, rq) // use gctx here
		if err != nil {
			return nil, err
		}
		defer func() {
			io.Copy(io.Discard, resp.Body)
			resp.Body.Close()
		}()

		if resp.StatusCode != http.StatusOK {
			b, _ := io.ReadAll(io.LimitReader(resp.Body, 8<<10))
			return nil, fmt.Errorf("get media url bad status: %s - %s", resp.Status, string(b))
		}

		d := &domain.DownloadLinkURL{}
		if err := json.NewDecoder(resp.Body).Decode(d); err != nil {
			return nil, fmt.Errorf("decode request: %w", err)
		}

		return s.DownloadMedia(ctx, d, fm) // DownloadMedia criará seu próprio dctx
	default:
		return nil, errors.New("not supported")
	}
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
