package domain

import "bytes"

type DownloadLinkURL struct {
	Url              string        `json:"url"`
	MimeType         string        `json:"mime_type"`
	Sha256           string        `json:"sha256"`
	FileSize         int           `json:"file_size"`
	Id               string        `json:"id"`
	MessagingProduct string        `json:"messaging_product"`
	UploadBuffer     *bytes.Buffer `json:"-"`
}
