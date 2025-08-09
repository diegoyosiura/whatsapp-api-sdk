package services

import (
	"context"

	"github.com/diegoyosiura/whatsapp-sdk-go/pkg/whatsapp/domain"
	"github.com/diegoyosiura/whatsapp-sdk-go/pkg/whatsapp/ports"
)

type PhoneService struct {
	api ports.PhoneAPI
}

func NewPhoneService(api ports.PhoneAPI) *PhoneService {
	return &PhoneService{api: api}
}

func (s *PhoneService) List(ctx context.Context) (*domain.PhoneList, error) {
	return s.api.List(ctx)
}

func (s *PhoneService) Get(ctx context.Context, phoneID string) (*domain.Phone, error) {
	return s.api.Get(ctx, phoneID)
}
