package ports

import "context"

// SecretsProvider retrieves non-rotating secrets required by the SDK, such as
// webhook verify token and app secret. SecretsProvider is intentionally generic
// to support multiple backends (env, file, redis, vault, etc.).
//
// Implementations SHOULD:
//   - Be safe for concurrent use by multiple goroutines.
//   - Support optional local caching to reduce backend latency/quotas.
//   - Mask values in logs; never print secret contents.
//
// Error semantics:
//   - Return an error when the key does not exist, is empty, or the backend fails.
//   - Wrap underlying backend failures with contextual information.
type SecretsProvider interface {
	// Get returns the secret value associated with the provided key.
	Get(ctx context.Context, key SecretKey) (string, error)
}
