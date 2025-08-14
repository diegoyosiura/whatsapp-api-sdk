package ports

import "bytes"

type SendMessage interface {
	Buffer() (*bytes.Buffer, error)
}
