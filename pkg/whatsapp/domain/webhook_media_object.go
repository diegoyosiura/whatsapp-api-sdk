package domain

type MediaObject struct {
	ID       string  `json:"id,omitempty"`
	MIMEType string  `json:"mime_type,omitempty"`
	SHA256   string  `json:"sha256,omitempty"`
	Caption  *string `json:"caption,omitempty"`
	Filename *string `json:"filename,omitempty"`
}
