package ports

import "context"

// TokenProvider supplies an access token for the Authorization header.
//
// Implementations SHOULD:
//   - Return a non-empty token quickly (prefer in-memory cache if backed by I/O).
//   - Be safe for concurrent use; coordinate refresh with single-flight if needed.
//   - Never log or expose token values in plaintext.
//
// Recommended error semantics:
//   - Return an error when the token is unavailable/expired/misconfigured.
//   - Wrap underlying backend failures (file/redis/vault) with context.
type TokenProvider interface {
	// Token returns the current bearer token. Implementations may internally
	// refresh/rotate as needed, but Token should aim to be fast and non-blocking.
	Token(ctx context.Context) (string, error)

	// Refresh forces token renewal when the transport detects 401 Unauthorized
	// or when proactive refresh is desired. Implementations may no-op if not supported.
	Refresh(ctx context.Context) error
}
