package domain

// CodeMethod represents how the verification code is delivered.
type CodeMethod string

const (
	CodeMethodSMS   CodeMethod = "SMS"
	CodeMethodVoice CodeMethod = "VOICE"
)

// RequestCodeParams is the pure domain payload for requesting a verification code.
// No HTTP concerns here.
type RequestCodeParams struct {
	CodeMethod CodeMethod `json:"code_method"`
	Locale     string     `json:"locale,omitempty"`
}

// VerifyCodeParams carries the code received by SMS/VOICE to confirm the phone.
type VerifyCodeParams struct {
	Code string `json:"code"`
}

// RegisterParams enables phone registration with an optional two-step PIN.
// "messaging_product" MUST be "whatsapp" (handled in transport layer).
type RegisterParams struct {
	Pin *string `json:"pin,omitempty"`
}

// TwoStepParams sets or updates the two-step verification code (PIN).
type TwoStepParams struct {
	Pin string `json:"pin"`
}

// ActionResult is the minimal common response for registration endpoints.
type ActionResult struct {
	Success bool `json:"success"`
}
