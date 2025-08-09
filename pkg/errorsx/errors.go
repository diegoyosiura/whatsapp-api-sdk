package errorsx

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// ErrNotConfigured indicates a required configuration value or dependency is missing.
var ErrNotConfigured = errors.New("not configured")

// HTTPError represents a non-2xx HTTP response from the transport layer.
// It captures enough context for actionable logs and debugging.
//
// NOTE: Body should be kept small (consider truncation at call-site) to avoid
// leaking large payloads into memory and logs.
type HTTPError struct {
	Method     string      // HTTP method used in the request
	URL        string      // Fully resolved URL that was requested
	StatusCode int         // HTTP status code
	Status     string      // HTTP status text
	Headers    http.Header // Response headers (may include rate-limit hints)
	Body       []byte      // Raw response body (truncated by caller if needed)
	FBTraceID  string      // Extracted fb-trace-id, when present
}

// Error implements the error interface.
func (e *HTTPError) Error() string {
	if e == nil {
		return "<nil>"
	}
	if e.FBTraceID != "" {
		return fmt.Sprintf("http %s %s: %d %s (fb-trace-id=%s)", e.Method, e.URL, e.StatusCode, e.Status, e.FBTraceID)
	}
	return fmt.Sprintf("http %s %s: %d %s", e.Method, e.URL, e.StatusCode, e.Status)
}

// IsRetryable reports whether the given error suggests a transient condition
// where a retry (with backoff) might succeed.
func IsRetryable(err error) bool {
	var he *HTTPError
	if errors.As(err, &he) {
		// Retry 429 and 5xx, except for 501/505 which usually indicate permanent issues.
		if he.StatusCode == http.StatusTooManyRequests {
			return true
		}
		if he.StatusCode >= 500 && he.StatusCode < 600 && he.StatusCode != http.StatusNotImplemented && he.StatusCode != http.StatusHTTPVersionNotSupported {
			return true
		}
	}
	return false
}

// NewHTTPErrorFromResponse constructs an HTTPError from an http.Response and body bytes.
// Callers should ensure body is already read (and optionally truncated) before calling.
func NewHTTPErrorFromResponse(resp *http.Response, body []byte) *HTTPError {
	if resp == nil {
		return &HTTPError{StatusCode: 0, Status: "<nil response>", Body: body}
	}
	fb := resp.Header.Get("x-fb-trace-id")
	return &HTTPError{
		Method:     resp.Request.Method,
		URL:        resp.Request.URL.String(),
		StatusCode: resp.StatusCode,
		Status:     resp.Status,
		Headers:    resp.Header.Clone(),
		Body:       body,
		FBTraceID:  fb,
	}
}

// GraphErrorDetail mirrors the canonical Facebook Graph error payload shape.
// See Meta Graph documentation for the authoritative schema.
type GraphErrorDetail struct {
	Message      string `json:"message"`
	Type         string `json:"type"`
	Code         int    `json:"code"`
	ErrorSubcode int    `json:"error_subcode,omitempty"`
	FBTraceID    string `json:"fbtrace_id,omitempty"`
}

// GraphError wraps an HTTPError with a decoded Graph error payload, when available.
// If the payload cannot be decoded, Raw will contain the original bytes for analysis.
type GraphError struct {
	HTTP   *HTTPError       // Underlying HTTP failure context
	Detail GraphErrorDetail // Decoded Graph error fields (may be zero-value)
	Raw    json.RawMessage  // Original raw body when decoding fails or for auditing
}

// Error implements the error interface.
func (e *GraphError) Error() string {
	if e == nil {
		return "<nil>"
	}
	if e.Detail.Message != "" {
		return fmt.Sprintf("graph error: %s (type=%s code=%d subcode=%d) — %s", e.Detail.Message, e.Detail.Type, e.Detail.Code, e.Detail.ErrorSubcode, e.HTTP.Error())
	}
	return fmt.Sprintf("graph error: undecoded — %s", e.HTTP.Error())
}

// Unwrap allows errors.Is/As to traverse to the underlying HTTPError.
func (e *GraphError) Unwrap() error { return e.HTTP }

// TryParseGraphError attempts to parse a Graph error payload from the given
// HTTP response body. On success, returns a *GraphError; otherwise returns nil.
func TryParseGraphError(resp *http.Response, body []byte) *GraphError {
	var envelope struct {
		Error GraphErrorDetail `json:"error"`
	}
	if err := json.Unmarshal(body, &envelope); err != nil {
		return &GraphError{HTTP: NewHTTPErrorFromResponse(resp, body), Raw: append([]byte(nil), body...)}
	}
	ge := &GraphError{
		HTTP:   NewHTTPErrorFromResponse(resp, body),
		Detail: envelope.Error,
		Raw:    append([]byte(nil), body...),
	}
	// Prefer fbtrace_id from payload if present.
	if ge.Detail.FBTraceID != "" {
		ge.HTTP.FBTraceID = ge.Detail.FBTraceID
	}
	return ge
}

// ValidationError represents client-side validation failures detected before
// issuing an HTTP request (e.g., missing required fields).
type ValidationError struct {
	Field  string // Name of the offending field (optional)
	Reason string // Human-friendly description of the problem
	Op     string // High-level operation name (e.g., "SendMessage")
}

// Error implements the error interface.
func (e *ValidationError) Error() string {
	if e == nil {
		return "<nil>"
	}
	if e.Field != "" {
		return fmt.Sprintf("validation error: op=%s field=%s: %s", e.Op, e.Field, e.Reason)
	}
	return fmt.Sprintf("validation error: op=%s: %s", e.Op, e.Reason)
}
