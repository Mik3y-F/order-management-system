package mpesa

import (
	"context"
	"fmt"

	"github.com/Mik3y-F/order-management-system/payments/internal/service"
	"github.com/Mik3y-F/order-management-system/pkg"
	"github.com/jwambugu/mpesa-golang-sdk"
)

const (
	MPESA_BUSINESS_SHORT_CODE = "MPESA_BUSINESS_SHORT_CODE" // #nosec G101 - This is an env variable name
	MPESA_PASSKEY             = "MPESA_PASSKEY"             // #nosec G101 - This is an env variable name
)

var _ service.PaymentsService = (*PaymentsService)(nil)

type PaymentsService struct {
	mpesa *Mpesa
}

func NewPaymentsService(mpesa *Mpesa) *PaymentsService {
	return &PaymentsService{
		mpesa: mpesa,
	}
}

func (s *PaymentsService) CheckPreconditions() {
	if s.mpesa == nil {
		panic("no Mpesa service provided")
	}
}

func (s *PaymentsService) ProcessPayment(ctx context.Context, payment *service.Payment) (*service.PaymentResponse, error) {
	s.CheckPreconditions()

	// stored in an environemnt variable for now:- assumption is that the system handles orders for a single business
	businessShortCode, err := pkg.StringToUint(pkg.MustGetEnv(MPESA_BUSINESS_SHORT_CODE))
	if err != nil {
		return nil, fmt.Errorf("failed to convert business short code to uint: %v", err)
	}

	passKey := pkg.MustGetEnv(MPESA_PASSKEY)

	stkPushRes, err := s.mpesa.app.STKPush(ctx, passKey, mpesa.STKPushRequest{
		BusinessShortCode: businessShortCode,
		TransactionType:   "CustomerBuyGoodsOnlines",
		Amount:            payment.Amount,
		PartyA:            payment.PhoneNumber,
		PartyB:            businessShortCode,
		PhoneNumber:       uint64(payment.PhoneNumber),
		CallBackURL:       payment.CallbackURL,
		AccountReference:  payment.Reference,
		TransactionDesc:   payment.Description,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to process payment: %v", MpesaErrorToInternalError(err))
	}

	return &service.PaymentResponse{
		MerchantRequestID: stkPushRes.MerchantRequestID,
		CheckoutRequestID: stkPushRes.CheckoutRequestID,
		ResponseCode:      stkPushRes.ResponseCode,
		CustomerMessage:   stkPushRes.CustomerMessage,
	}, nil

}
