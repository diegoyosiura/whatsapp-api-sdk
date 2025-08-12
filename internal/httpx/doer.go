package httpx

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"syscall"
	"time"
)

// RetryPolicy reports whether an HTTP status code is retryable.
type RetryPolicy func(status int) bool

// Options configures the Doer.
type Options struct {
	// Transport defines the underlying RoundTripper. If nil, http.DefaultTransport is used.
	Transport http.RoundTripper
	// MaxRetries defines how many retry attempts (in addition to the initial try)
	// will be performed for retryable responses. Default 3 when zero.
	MaxRetries int
	// BaseBackoff is the initial backoff duration. Default 200ms when zero.
	BaseBackoff time.Duration
	// MaxBackoff caps the backoff. Default 3s when zero.
	MaxBackoff time.Duration
	// RetryPolicy determines which statuses are retryable. If nil, DefaultRetryPolicy is used.
	RetryPolicy RetryPolicy
}

// Doer implements a context-aware HTTP executor with basic retry and jitter.
type Doer struct {
	client      *http.Client
	maxRetries  int
	baseBackoff time.Duration
	maxBackoff  time.Duration
	shouldRetry RetryPolicy
}

// New creates a new Doer with the provided options.
func New(opts Options) *Doer {
	tr := opts.Transport
	if tr == nil {
		tr = http.DefaultTransport
	}
	maxRetries := opts.MaxRetries
	if maxRetries == 0 {
		maxRetries = 3
	}
	bb := opts.BaseBackoff
	if bb <= 0 {
		bb = 200 * time.Millisecond
	}
	mb := opts.MaxBackoff
	if mb <= 0 {
		mb = 3 * time.Second
	}
	pol := opts.RetryPolicy
	if pol == nil {
		pol = DefaultRetryPolicy
	}
	return &Doer{
		client:      &http.Client{Transport: tr},
		maxRetries:  maxRetries,
		baseBackoff: bb,
		maxBackoff:  mb,
		shouldRetry: pol,
	}
}

// DefaultRetryPolicy retries 429 and 5xx except for 501/505.
func DefaultRetryPolicy(status int) bool {
	if status == http.StatusTooManyRequests {
		return true
	}
	if status >= 500 && status < 600 && status != http.StatusNotImplemented && status != http.StatusHTTPVersionNotSupported {
		return true
	}
	return false
}

// Do executes req with retry for retryable statuses. It honors ctx for timeout/cancel.
// IMPORTANT: req.Body must be rewindable for retries. The caller should set
// GetBody on the request (Go 1.19+) or provide a fresh request per attempt.
func (d *Doer) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	if req == nil {
		return nil, errors.New("httpx: nil request")
	}
	// Ensure the request uses the provided context for deadlines/cancellations.
	req = req.WithContext(ctx)

	var lastErr error
	var resp *http.Response

	for attempt := 0; attempt <= d.maxRetries; attempt++ {
		// Rewind body if necessary and possible.
		if attempt > 0 {
			if req.GetBody != nil {
				body, err := req.GetBody()
				if err != nil {
					return nil, fmt.Errorf("httpx: rewind body: %w", err)
				}
				req.Body = body
			} else if req.Body != nil {
				// If body cannot be rewound, we cannot retry safely.
				return nil, errors.New("httpx: request body is not rewindable; set Request.GetBody for retries")
			}
		}

		resp, lastErr = d.client.Do(req)
		if lastErr != nil {
			// Network or context error – do not blindly retry on permanent failures.
			if isTempOrTimeout(lastErr) && attempt < d.maxRetries {
				d.sleep(backoff(attempt, d.baseBackoff, d.maxBackoff))
				continue
			}
			return nil, lastErr
		}

		// If non-nil response, decide on retry based on status code.
		if !d.shouldRetry(resp.StatusCode) || attempt == d.maxRetries {
			return resp, nil
		}
		// Drain and close the body before retrying to reuse connections.
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
		d.sleep(backoff(attempt, d.baseBackoff, d.maxBackoff))
	}

	// Should not reach here.
	return resp, lastErr
}

func (d *Doer) sleep(dur time.Duration) {
	// Sleep respecting context is handled at a higher level; here we simply wait.
	time.Sleep(dur)
}

func isTempOrTimeout(err error) bool {
	if err == nil {
		return false
	}

	// Permanent: caller canceled or listener/conn closed.
	if errors.Is(err, context.Canceled) || errors.Is(err, net.ErrClosed) {
		return false
	}

	// Generic network timeout.
	var ne net.Error
	if errors.As(err, &ne) && ne.Timeout() {
		return true
	}

	// DNS specific signals.
	var dnserr *net.DNSError
	if errors.As(err, &dnserr) {
		if dnserr.IsTimeout || dnserr.IsTemporary {
			return true
		}
	}

	// Common transient syscall errors seen on sockets.
	switch {
	case errors.Is(err, syscall.ECONNRESET),
		errors.Is(err, syscall.ECONNABORTED),
		errors.Is(err, syscall.EPIPE),
		errors.Is(err, syscall.ETIMEDOUT),
		errors.Is(err, syscall.EHOSTUNREACH),
		errors.Is(err, syscall.ENETUNREACH):
		return true
	}

	// Often transient when reading HTTP bodies.
	if errors.Is(err, io.ErrUnexpectedEOF) {
		return true
	}

	return false
}

func backoff(attempt int, base, max time.Duration) time.Duration {
	// Exponential backoff: base * 2^attempt, with jitter up to 50%.
	m := 1 << attempt
	d := time.Duration(m) * base
	if d > max {
		d = max
	}
	// Jitter in [0, d/2]
	jit := time.Duration(rand.Int63n(int64(d / 2)))
	return d/2 + jit
}
