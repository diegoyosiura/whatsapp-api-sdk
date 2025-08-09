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

// MessagesEndpoint returns the full URL for POST /{Version}/{Phone-Number-ID}/messages
func MessagesEndpoint(base, version, phoneNumberID string) string {
	return buildURL(base, version, phoneNumberID, "messages")
}

// PhoneNumbersListEndpoint returns GET /{Version}/{WABA-ID}/phone_numbers
func PhoneNumbersListEndpoint(base, version, wabaID string) string {
	return buildURL(base, version, wabaID, "phone_numbers")
}

// PhoneNumberGetEndpoint returns GET /{Version}/{Phone-Number-ID}
func PhoneNumberGetEndpoint(base, version, phoneNumberID string) string {
	return buildURL(base, version, phoneNumberID)
}

// RegisterEndpoint Registration endpoints
func RegisterEndpoint(base, version, phoneNumberID string) string {
	return buildURL(base, version, phoneNumberID, "register")
}

// DeregisterEndpoint DeregisterEndpoint
func DeregisterEndpoint(base, version, phoneNumberID string) string {
	return buildURL(base, version, phoneNumberID, "deregister")
}

// RequestCodeEndpoint RequestCodeEndpoint
func RequestCodeEndpoint(base, version, phoneNumberID string) string {
	return buildURL(base, version, phoneNumberID, "request_code")
}

// VerifyCodeEndpoint VerifyCodeEndpoint
func VerifyCodeEndpoint(base, version, phoneNumberID string) string {
	return buildURL(base, version, phoneNumberID, "verify_code")
}

// TwoFactorEndpoint Two-factor (2FA) uses POST /{Version}/{Phone-Number-ID} with fields.
func TwoFactorEndpoint(base, version, phoneNumberID string) string {
	return buildURL(base, version, phoneNumberID)
}
