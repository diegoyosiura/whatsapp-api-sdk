package domain

type InteractiveObject struct {
	Type string `json:"type"`
	// Add subtypes as needed, e.g., button replies or list replies
	// For example: "button" replies emit a nested structure with id/title.
}
