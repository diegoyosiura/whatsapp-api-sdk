package domain

import "bytes"

type UploadRequest struct {
	Url              string        `json:"url"`
	MessagingProduct string        `json:"messaging_product"`
	UploadBuffer     *bytes.Buffer `json:"-"`
}
