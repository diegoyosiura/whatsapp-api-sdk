package utils

import (
	"fmt"
	"net/http"
	"net/textproto"
	"os"
	"path/filepath"
	"slices"
)

func GetMimeTypeHeader(f *os.File) textproto.MIMEHeader {
	ext := filepath.Ext(f.Name())
	mimeType := ""
	if slices.Contains([]string{".xlsx", ".docx", ".pptx"}, ext) {
		mimeType = map[string]string{
			".docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
			".xlsx": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
			".pptx": "application/vnd.openxmlformats-officedocument.presentationml.presentation",
		}[ext]
	} else {
		head := make([]byte, 512)
		n, _ := f.Read(head)
		mimeType = http.DetectContentType(head[:n])
	}

	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file"; filename="%s"`, filepath.Base(f.Name())))
	h.Set("Content-Type", mimeType)

	return h
}
