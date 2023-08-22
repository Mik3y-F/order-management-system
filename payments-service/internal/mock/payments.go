package mock

import (
	"context"

	"github.com/Mik3y-F/order-management-system/payments/internal/service"
)

var _ service.PaymentsService = (*PaymentsService)(nil)

type PaymentsService struct {
	ProcessPaymentFunc func(ctx context.Context, p *service.Payment) (*service.PaymentResponse, error)
}

func (m *PaymentsService) ProcessPayment(ctx context.Context, p *service.Payment) (*service.PaymentResponse, error) {
	return m.ProcessPaymentFunc(ctx, p)
}
