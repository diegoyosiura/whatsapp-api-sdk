// Package httpx implements an internal HTTP client with retry and backoff.
//
// Design notes:
//   - The SDK applies per-call timeouts via context; httpx does not set a
//     global http.Client.Timeout to preserve streaming behavior.
//   - Retries only occur for retryable responses (429/5xx) or network
//     timeouts marked by net.Error. Request bodies must be rewindable.
//   - No external dependencies; callers can wrap with their own telemetry.
package httpx
