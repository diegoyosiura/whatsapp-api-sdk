package ports

// SecretKey represents a logical key for non-rotating secrets that the SDK needs.
// Typical keys include webhook verify token and app secret (HMAC key).
// This typed alias increases readability and reduces stringly-typed mistakes.
type SecretKey string

const (
	// VerifyTokenKey is used to retrieve the webhook verify token.
	VerifyTokenKey SecretKey = "verify_token"
	// AppSecretKey is used to retrieve the app secret for HMAC signature validation.
	AppSecretKey SecretKey = "app_secret"
)
