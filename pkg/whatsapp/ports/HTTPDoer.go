package ports

import (
	"context"
	"net/http"
)

// HTTPDoer abstracts HTTP transport execution for the Meta Graph API.
//
// Implementations SHOULD:
//   - Honor context deadlines and cancellations.
//   - Be safe for concurrent use by multiple goroutines.
//   - Preserve headers set by the caller (Authorization, User-Agent, etc.).
//   - Optionally implement retry policies for 429/5xx (with backoff/jitter).
//
// Implementations MUST NOT:
//   - Encode/decode domain models (that is the caller's responsibility).
//   - Swallow context errors (return context.DeadlineExceeded / context.Canceled as-is).
//
// The caller is responsible for closing the response body when non-nil.
type HTTPDoer interface {
	// Do executes the given HTTP request and returns the raw HTTP response or an error.
	// The request is assumed to be fully constructed (URL, method, headers, body).
	Do(ctx context.Context, req *http.Request) (*http.Response, error)
}
