package ports

import (
	"context"
	"fmt"

	"github.com/diegoyosiura/whatsapp-sdk-go/pkg/whatsapp/ports"
)

// FakeSecretsProvider is an in-memory SecretsProvider that resolves secrets
// from a simple map keyed by ports.SecretKey.
type FakeSecretsProvider struct {
	Secrets map[ports.SecretKey]string
	Err     error
}

// Get returns the value in the map or an error if Err is set or the key is missing.
func (f *FakeSecretsProvider) Get(ctx context.Context, key ports.SecretKey) (string, error) {
	if f.Err != nil {
		return "", f.Err
	}
	if val, ok := f.Secrets[key]; ok {
		return val, nil
	}
	return "", fmt.Errorf("secret %q not found", key)
}
