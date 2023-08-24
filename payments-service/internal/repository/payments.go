package repository

import "context"

type PaymentStatus string

const (
	PaymentStatusPending PaymentStatus = "pending"
	PaymentStatusPaid    PaymentStatus = "paid"
	PaymentStatusFailed  PaymentStatus = "failed"
)

type Payment struct {
	Id          string
	Amount      uint
	Status      PaymentStatus
	OrderID     string
	Phone       string
	Reference   string
	Description string
	CreatedAt   string
	UpdatedAt   string
}

func (p *Payment) Validate() error {
	return nil
}

type PaymentsRepository interface {
	CreatePayment(ctx context.Context, payment *Payment) (string, error)
	GetPaymentByID(ctx context.Context, paymentID string) (*Payment, error)
	UpdatePaymentStatus(ctx context.Context, paymentID string, status PaymentStatus) error
}
