package service

import "context"

type Payment struct {
	Id          string `json:"id"`
	OrderId     string `json:"orderId"`
	PhoneNumber uint   `json:"phoneNumber"`
	Amount      uint   `json:"amount"`
	Reference   string `json:"reference"`
	Description string `json:"description"`
	CallbackURL string `json:"callbackUrl"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

type PaymentResponse struct {
	CheckoutRequestID string `json:"CheckoutRequestID"`
	CustomerMessage   string `json:"CustomerMessage"`
	MerchantRequestID string `json:"MerchantRequestID"`
	ResponseCode      string `json:"ResponseCode"`
	ResponseMessage   string `json:"ResponseMessage"`
}

type PaymentsService interface {
	ProcessPayment(ctx context.Context, payment *Payment) (*PaymentResponse, error)
	// CheckPayment(ctx context.Context, id string) (*Payment, error)
}
