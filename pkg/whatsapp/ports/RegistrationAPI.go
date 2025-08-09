package ports

import (
	"context"

	"github.com/diegoyosiura/whatsapp-sdk-go/pkg/whatsapp/domain"
)

// RegistrationAPI descreve as operações de registro de número.
type RegistrationAPI interface {
	RequestCode(ctx context.Context, p domain.RequestCodeParams) (*domain.ActionResult, error)
	VerifyCode(ctx context.Context, p domain.VerifyCodeParams) (*domain.ActionResult, error)
	Register(ctx context.Context, p domain.RegisterParams) (*domain.ActionResult, error)
	Deregister(ctx context.Context) (*domain.ActionResult, error)
	SetTwoStep(ctx context.Context, p domain.TwoStepParams) (*domain.ActionResult, error)
}
