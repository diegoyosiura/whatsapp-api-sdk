package whatsapp

import (
	"errors"
	"fmt"
	"time"

	"github.com/diegoyosiura/whatsapp-sdk-go/pkg/errorsx"
	"github.com/diegoyosiura/whatsapp-sdk-go/pkg/whatsapp/ports"
)

// Options configures the SDK Client. Callers should prefer passing explicit
// Options instead of relying solely on environment variables, to make behavior
// reproducible and testable.
type Options struct {
	// Graph API version in the form vMAJOR.MINOR, e.g., "v20.0".
	Version string
	// Business account and phone number identifiers.
	WABAID        string
	PhoneNumberID string

	// Transport and providers (required).
	HTTPDoer        ports.HTTPDoer
	TokenProvider   ports.TokenProvider
	SecretsProvider ports.SecretsProvider // optional unless using features that need it

	// Optional base URL override (for testing/mocking). Empty means use default Graph API base.
	BaseURL string

	// Network and resilience settings.
	Timeout   time.Duration // per-request timeout
	RetryMax  int           // max retries for retryable statuses
	UserAgent string        // appended to default UA if non-empty
}

// Validate checks that Options contain a minimal viable configuration.
func (o *Options) Validate() error {
	if o == nil {
		return errors.New("options: nil")
	}
	if o.Version == "" {
		return &errorsx.ValidationError{Op: "ClientInit", Field: "Version", Reason: "empty"}
	}
	if len(o.Version) < 4 || o.Version[0] != 'v' {
		return &errorsx.ValidationError{Op: "ClientInit", Field: "Version", Reason: "must be like v20.0"}
	}
	if o.WABAID == "" {
		return &errorsx.ValidationError{Op: "ClientInit", Field: "WABAID", Reason: "empty"}
	}
	if o.PhoneNumberID == "" {
		return &errorsx.ValidationError{Op: "ClientInit", Field: "PhoneNumberID", Reason: "empty"}
	}
	if o.TokenProvider == nil {
		return &errorsx.ValidationError{Op: "ClientInit", Field: "TokenProvider", Reason: "nil"}
	}
	return nil
}

// withDefaults fills zero-values with sensible defaults.
func (o *Options) withDefaults() Options {
	cpy := *o
	if cpy.Timeout == 0 {
		cpy.Timeout = 10 * time.Second
	}
	if cpy.RetryMax == 0 {
		cpy.RetryMax = 3
	}
	return cpy
}

func (o Options) String() string {
	// Do not print secrets; only structural info.
	return fmt.Sprintf("Options{Version=%s WABAID=%s PhoneNumberID=%s BaseURL=%s Timeout=%s RetryMax=%d UA+=%t}",
		o.Version, mask(o.WABAID), mask(o.PhoneNumberID), o.BaseURL, o.Timeout, o.RetryMax, o.UserAgent != "")
}

func mask(s string) string {
	if len(s) <= 4 {
		return "****"
	}
	return s[:2] + "****" + s[len(s)-2:]
}
