package graph

import (
	"net/url"
	"path"
)

const DefaultBaseURL = "https://graph.facebook.com"

// buildURL joins base, version and parts into a clean URL string.
func buildURL(base, version string, parts ...string) string {
	u, _ := url.Parse(base)
	p := append([]string{version}, parts...)
	u.Path = path.Join(append([]string{u.Path}, p...)...)
	return u.String()
}

// MessagesEndpoint returns the full URL for POST /{Version}/{Phone-Number-ID}/messages.
func MessagesEndpoint(base, version, phoneNumberID string) string {
	return buildURL(base, version, phoneNumberID, "messages")
}

// PhoneNumbersListEndpoint returns the full URL for GET /{Version}/{WABA-ID}/phone_numbers.
func PhoneNumbersListEndpoint(base, version, wabaID string) string {
	return buildURL(base, version, wabaID, "phone_numbers")
}

// PhoneNumberGetEndpoint returns the full URL for GET /{Version}/{Phone-Number-ID}.
func PhoneNumberGetEndpoint(base, version, phoneNumberID string) string {
	return buildURL(base, version, phoneNumberID)
}

// RegisterEndpoint returns the full URL for POST /{Version}/{Phone-Number-ID}/register.
func RegisterEndpoint(base, version, phoneNumberID string) string {
	return buildURL(base, version, phoneNumberID, "register")
}

// DeregisterEndpoint returns the full URL for POST /{Version}/{Phone-Number-ID}/deregister.
func DeregisterEndpoint(base, version, phoneNumberID string) string {
	return buildURL(base, version, phoneNumberID, "deregister")
}

// RequestCodeEndpoint returns the full URL for POST /{Version}/{Phone-Number-ID}/request_code.
func RequestCodeEndpoint(base, version, phoneNumberID string) string {
	return buildURL(base, version, phoneNumberID, "request_code")
}

// VerifyCodeEndpoint returns the full URL for POST /{Version}/{Phone-Number-ID}/verify_code.
func VerifyCodeEndpoint(base, version, phoneNumberID string) string {
	return buildURL(base, version, phoneNumberID, "verify_code")
}

// TwoFactorEndpoint returns the full URL for POST /{Version}/{Phone-Number-ID} used for two-factor (2FA) requests.
func TwoFactorEndpoint(base, version, phoneNumberID string) string {
	return buildURL(base, version, phoneNumberID)
}
func RequestMediaLink(base, version, mediaID string) string {
	return buildURL(base, version, mediaID)
}
