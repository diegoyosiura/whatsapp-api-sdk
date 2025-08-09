package domain

// QualityRating is the health of the phone number's messaging quality.
type QualityRating string

const (
	QualityGreen  QualityRating = "GREEN"
	QualityYellow QualityRating = "YELLOW"
	QualityRed    QualityRating = "RED"
)

// Phone represents a phone number entity returned by the Graph API.
type Phone struct {
	ID                 string        `json:"id"`
	DisplayPhoneNumber string        `json:"display_phone_number"`
	VerifiedName       string        `json:"verified_name"`
	QualityRating      QualityRating `json:"quality_rating,omitempty"`
	IsOfficialBusiness *bool         `json:"is_official_business_account,omitempty"`
	AccountMode        string        `json:"account_mode,omitempty"` // e.g. "SANDBOX", "LIVE"
	// Add more fields as you need from the collection (country, certificate, etc.)
}
