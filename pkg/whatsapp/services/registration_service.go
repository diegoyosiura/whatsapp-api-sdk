package services

import (
	"context"

	"github.com/diegoyosiura/whatsapp-sdk-go/pkg/whatsapp/domain"
	"github.com/diegoyosiura/whatsapp-sdk-go/pkg/whatsapp/ports"
)

type RegistrationService struct {
	api ports.RegistrationAPI
}

func NewRegistrationService(api ports.RegistrationAPI) *RegistrationService {
	return &RegistrationService{api: api}
}

func (s *RegistrationService) RequestCode(ctx context.Context, p domain.RequestCodeParams) (*domain.ActionResult, error) {
	return s.api.RequestCode(ctx, p)
}
func (s *RegistrationService) VerifyCode(ctx context.Context, p domain.VerifyCodeParams) (*domain.ActionResult, error) {
	return s.api.VerifyCode(ctx, p)
}
func (s *RegistrationService) Register(ctx context.Context, p domain.RegisterParams) (*domain.ActionResult, error) {
	return s.api.Register(ctx, p)
}
func (s *RegistrationService) Deregister(ctx context.Context) (*domain.ActionResult, error) {
	return s.api.Deregister(ctx)
}
func (s *RegistrationService) SetTwoStep(ctx context.Context, p domain.TwoStepParams) (*domain.ActionResult, error) {
	return s.api.SetTwoStep(ctx, p)
}
func (s *RegistrationService) API() ports.RegistrationAPI {
	return s.api
}
