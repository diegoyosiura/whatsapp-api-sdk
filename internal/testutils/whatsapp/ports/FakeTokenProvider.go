package ports

import "context"

// FakeTokenProvider is an in-memory TokenProvider for unit testing. It returns
// a static token and errors as configured in the struct fields.
type FakeTokenProvider struct {
	TokenValue string
	Err        error
	RefreshErr error
}

// Token returns TokenValue or Err when set.
func (f *FakeTokenProvider) Token(ctx context.Context) (string, error) { return f.TokenValue, f.Err }

// Refresh returns RefreshErr when set.
func (f *FakeTokenProvider) Refresh(ctx context.Context) error { return f.RefreshErr }
