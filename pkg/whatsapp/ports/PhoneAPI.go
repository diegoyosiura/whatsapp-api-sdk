package ports

import (
	"context"

	"github.com/diegoyosiura/whatsapp-sdk-go/pkg/whatsapp/domain"
)

type PhoneAPI interface {
	List(ctx context.Context) (*domain.PhoneList, error)
	Get(ctx context.Context, phoneID string) (*domain.Phone, error)
}
