package handlers

import (
	"context"
	"fmt"

	pb "github.com/Mik3y-F/order-management-system/orders/api/generated"
)

func (s *GRPCServer) ProcessCheckout(
	ctx context.Context, in *pb.ProcessCheckoutRequest) (*pb.ProcessCheckoutResponse, error) {

	o, err := s.CheckoutService.ProcessCheckout(ctx, in.GetOrderId())
	if err != nil {
		return nil, Error(fmt.Errorf("failed to process checkout: %w", err))
	}

	return &pb.ProcessCheckoutResponse{
		OrderId:    o.Id,
		CustomerId: o.CustomerId,
		Status:     getGRPCOrderStatus(o.OrderStatus),
	}, nil
}
