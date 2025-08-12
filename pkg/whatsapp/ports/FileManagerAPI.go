package ports

import "context"

type FileManagerAPI interface {
	SetData(ctx context.Context, data []byte) error
	Save(ctx context.Context, fileName string) error
	Open(ctx context.Context, fileName string) ([]byte, error)
}
