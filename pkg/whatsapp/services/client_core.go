package services

import (
	"context"
	"net/http"

	"github.com/diegoyosiura/whatsapp-sdk-go/pkg/whatsapp/ports"
)

// MessagesService provides high-level operations for sending messages.
// It uses the transport/graph adapter to build HTTP requests and relies on the
// Client's configured HTTPDoer and TokenProvider to execute them.
type MessagesService struct {
	c clientCore
}

// clientCore is the minimal facade the service needs; *whatsapp.Client satisfies it.
type clientCore interface {
	BaseURL() string
	Version() string
	PhoneNumberID() string
	TokenProvider() ports.TokenProvider
	Do(ctx context.Context, req *http.Request) (*http.Response, error)
}
