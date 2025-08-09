package whatsapp

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/diegoyosiura/whatsapp-sdk-go/internal/httpx"
	"github.com/diegoyosiura/whatsapp-sdk-go/pkg/whatsapp/ports"
	"github.com/diegoyosiura/whatsapp-sdk-go/pkg/whatsapp/services"
	"github.com/diegoyosiura/whatsapp-sdk-go/pkg/whatsapp/transport/graph"
)

// Client is the public facade of the SDK. It aggregates services and holds
// shared configuration (version, IDs, providers, transport).
type Client struct {
	version       string
	wabaID        string
	phoneNumberID string

	httpDoer      ports.HTTPDoer
	tokenProvider ports.TokenProvider
	secrets       ports.SecretsProvider

	Messages     *services.MessagesService
	Phone        *services.PhoneService
	Registration *services.RegistrationService
	Webhook      *services.WebhookService

	baseURL  string
	timeout  time.Duration
	retryMax int
	uaExtra  string
}

// NewClient validates options, applies defaults and returns a ready-to-use Client.
// If no HTTPDoer is provided, a default httpx.Doer is constructed using RetryMax
// from Options and the default RoundTripper.
func NewClient(o Options) (*Client, error) {
	if err := o.Validate(); err != nil {
		return nil, err
	}
	o = o.withDefaults()

	var doer ports.HTTPDoer = o.HTTPDoer
	if doer == nil {
		doer = httpx.New(httpx.Options{MaxRetries: o.RetryMax})
	}

	phoneAPI := graph.NewPhoneAPI(doer, o.TokenProvider, o.Version, o.WABAID)
	regAPI := graph.NewRegistrationAPI(doer, o.TokenProvider, o.Version, o.PhoneNumberID)

	c := &Client{
		Phone:         services.NewPhoneService(phoneAPI),
		Registration:  services.NewRegistrationService(regAPI),
		Webhook:       services.NewWebhookService(o.SecretsProvider),
		version:       o.Version,
		wabaID:        o.WABAID,
		phoneNumberID: o.PhoneNumberID,
		httpDoer:      doer,
		tokenProvider: o.TokenProvider,
		secrets:       o.SecretsProvider,
		baseURL:       o.BaseURL,
		timeout:       o.Timeout,
		retryMax:      o.RetryMax,
		uaExtra:       o.UserAgent,
	}
	c.Messages = services.NewMessagesService(c)
	return c, nil
}

// do constructs and executes an HTTP request using the configured HTTPDoer.
// It injects Authorization and User-Agent headers, applies timeout, and returns
// the raw *http.Response for the caller to decode.
func (c *Client) do(ctx context.Context, req *http.Request) (*http.Response, error) {
	// Apply per-request timeout via context.
	var cancel context.CancelFunc
	if c.timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, c.timeout)
		defer cancel()
	}

	// Inject Authorization header.
	token, err := c.tokenProvider.Token(ctx)
	if err != nil {
		return nil, fmt.Errorf("get token: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	// User-Agent (append extra if provided).
	ua := "ampere-whatsapp-sdk-go"
	if c.uaExtra != "" {
		ua = ua + " " + c.uaExtra
	}
	req.Header.Set("User-Agent", ua)

	// Execute via injected transport.
	resp, err := c.httpDoer.Do(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Narrow getters used by services without leaking internal fields.

func (c *Client) Version() string                    { return c.version }
func (c *Client) WABAID() string                     { return c.wabaID }
func (c *Client) PhoneNumberID() string              { return c.phoneNumberID }
func (c *Client) BaseURL() string                    { return c.baseURL }
func (c *Client) TokenProvider() ports.TokenProvider { return c.tokenProvider }
func (c *Client) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	return c.do(ctx, req)
}
func (c *Client) RegistrationService() *services.RegistrationService {
	return services.NewRegistrationService(
		graph.NewRegistrationAPI(c.httpDoer, c.tokenProvider, c.version, c.phoneNumberID),
	)
}
