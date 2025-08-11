package graph

import "testing"

func TestBuildURLAndEndpoints(t *testing.T) {
	base := "https://graph.example.com"
	version := "v1"
	phoneID := "123"
	wabaID := "waba"

	if got := buildURL(base, version, "a", "b"); got != "https://graph.example.com/v1/a/b" {
		t.Fatalf("buildURL unexpected: %s", got)
	}

	cases := []struct {
		name string
		got  string
		want string
	}{
		{"messages", MessagesEndpoint(base, version, phoneID), "https://graph.example.com/v1/123/messages"},
		{"phoneNumbersList", PhoneNumbersListEndpoint(base, version, wabaID), "https://graph.example.com/v1/waba/phone_numbers"},
		{"phoneNumberGet", PhoneNumberGetEndpoint(base, version, phoneID), "https://graph.example.com/v1/123"},
		{"register", RegisterEndpoint(base, version, phoneID), "https://graph.example.com/v1/123/register"},
		{"deregister", DeregisterEndpoint(base, version, phoneID), "https://graph.example.com/v1/123/deregister"},
		{"requestCode", RequestCodeEndpoint(base, version, phoneID), "https://graph.example.com/v1/123/request_code"},
		{"verifyCode", VerifyCodeEndpoint(base, version, phoneID), "https://graph.example.com/v1/123/verify_code"},
		{"twoFactor", TwoFactorEndpoint(base, version, phoneID), "https://graph.example.com/v1/123"},
	}

	for _, tt := range cases {
		if tt.got != tt.want {
			t.Errorf("%s: got %s, want %s", tt.name, tt.got, tt.want)
		}
	}
}
