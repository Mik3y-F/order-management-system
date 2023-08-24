package grpc

import (
	"context"

	pb "github.com/Mik3y-F/order-management-system/payments/api/generated"
	"github.com/Mik3y-F/order-management-system/payments/internal/service"
)

func (s *GRPCServer) ProcessPayment(ctx context.Context, in *pb.MpesaPaymentRequest) (*pb.MpesaPaymentResponse, error) {

	callBackURL := "https://order-management-system.herokuapp.com/payments/callback"

	p, err := s.PaymentsService.ProcessPayment(ctx, &service.Payment{
		Amount:      uint(in.GetAmount()),
		PhoneNumber: uint(in.GetPhoneNumber()),
		CallbackURL: callBackURL,
		Reference:   in.GetReference(),
		Description: in.GetDescription(),
	})
	if err != nil {
		return nil, Error(err)
	}

	return &pb.MpesaPaymentResponse{
		CheckoutRequestId: p.CheckoutRequestID,
		CustomerMessage:   p.CustomerMessage,
		ResponseCode:      p.ResponseCode,
		MerchantRequestId: p.MerchantRequestID,
	}, nil
}
