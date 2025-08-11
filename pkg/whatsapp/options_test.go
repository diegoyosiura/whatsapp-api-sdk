package whatsapp

import (
	"context"
	"strings"
	"testing"
	"time"
)

type stubTokenProvider struct{}

func (stubTokenProvider) Token(ctx context.Context) (string, error) { return "token", nil }
func (stubTokenProvider) Refresh(ctx context.Context) error         { return nil }

func validOpts() Options {
	return Options{
		Version:       "v1.0",
		WABAID:        "12345678",
		PhoneNumberID: "87654321",
		TokenProvider: stubTokenProvider{},
	}
}

func TestOptionsValidate(t *testing.T) {
	o := validOpts()
	if err := o.Validate(); err != nil {
		t.Fatalf("valid options failed: %v", err)
	}
	cases := []struct {
		name string
		mod  func(o *Options)
	}{
		{"nil", func(o *Options) { *o = Options{} }},
		{"no version", func(o *Options) { o.Version = "" }},
		{"bad version", func(o *Options) { o.Version = "20.0" }},
		{"no waba", func(o *Options) { o.WABAID = "" }},
		{"no phone", func(o *Options) { o.PhoneNumberID = "" }},
		{"no token", func(o *Options) { o.TokenProvider = nil }},
	}
	for _, tc := range cases {
		o := validOpts()
		tc.mod(&o)
		if err := o.Validate(); err == nil {
			t.Fatalf("%s: expected error", tc.name)
		}
	}
}

func TestOptionsWithDefaults(t *testing.T) {
	o := validOpts()
	d := o.withDefaults()
	if d.Timeout != 10*time.Second || d.RetryMax != 3 {
		t.Fatalf("unexpected defaults: %+v", d)
	}
	o.Timeout = time.Second
	o.RetryMax = 5
	d = o.withDefaults()
	if d.Timeout != o.Timeout || d.RetryMax != o.RetryMax {
		t.Fatalf("non-zero values should be kept")
	}
}

func TestOptionsStringMasks(t *testing.T) {
	o := validOpts()
	o.BaseURL = "https://example"
	o.Timeout = time.Second
	o.RetryMax = 5
	o.UserAgent = "ua"
	s := o.String()
	if !strings.Contains(s, "Options{Version=v1.0") {
		t.Fatalf("unexpected string: %s", s)
	}
	if !strings.Contains(s, "WABAID=12****78") || !strings.Contains(s, "PhoneNumberID=87****21") {
		t.Fatalf("ids not masked: %s", s)
	}
}
