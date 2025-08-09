package ports

import (
	"context"
	"net/http"
)

type FakeHTTPDoer struct {
	Fn func(ctx context.Context, req *http.Request) (*http.Response, error)
}

// Do executes the function Fn if set; otherwise returns a nil response and nil error.
func (f *FakeHTTPDoer) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	if f.Fn != nil {
		return f.Fn(ctx, req)
	}
	return nil, nil
}
