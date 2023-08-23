package handlers

import (
	"context"
	"fmt"

	pb "github.com/Mik3y-F/order-management-system/orders/api/generated"
	"github.com/Mik3y-F/order-management-system/orders/internal/service"
)

func (s *GRPCServer) CreateOrder(ctx context.Context, in *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {

	var items []*service.OrderItem
	for _, item := range in.GetOrderItems() {
		items = append(items, &service.OrderItem{
			ProductId: item.GetProductId(),
			Quantity:  uint(item.GetQuantity()),
		})
	}

	p, err := s.OrderService.CreateOrder(ctx, &service.Order{
		CustomerId: in.GetCustomerId(),
		Items:      items,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	return &pb.CreateOrderResponse{
		Id: p.Id,
	}, nil
}

func (s *GRPCServer) GetOrder(ctx context.Context, in *pb.GetOrderRequest) (*pb.GetOrderResponse, error) {

	order, err := s.OrderService.GetOrder(ctx, in.GetId())
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	var orderItems []*pb.OrderItem
	for _, item := range order.Items {
		orderItems = append(orderItems, &pb.OrderItem{
			ProductId: item.ProductId,
			Quantity:  uint32(item.Quantity),
		})
	}

	return &pb.GetOrderResponse{
		Id:         order.Id,
		CustomerId: order.CustomerId,
		OrderItems: orderItems,
		Status:     getGRPCOrderStatus(order.OrderStatus),
	}, nil
}

func (s *GRPCServer) ListOrders(ctx context.Context, in *pb.ListOrdersRequest) (*pb.ListOrdersResponse, error) {

	orders, err := s.OrderService.ListOrders(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list orders: %w", err)
	}

	var responseOrders []*pb.Order
	for _, p := range orders {
		var orderItems []*pb.OrderItem
		for _, item := range p.Items {
			orderItems = append(orderItems, &pb.OrderItem{
				ProductId: item.ProductId,
				Quantity:  uint32(item.Quantity),
			})
		}

		responseOrders = append(responseOrders, &pb.Order{
			Id:         p.Id,
			CustomerId: p.CustomerId,
			OrderItems: orderItems,
			Status:     getGRPCOrderStatus(p.OrderStatus),
		})
	}

	return &pb.ListOrdersResponse{
		Orders: responseOrders,
	}, nil
}

func getGRPCOrderStatus(status service.OrderStatus) pb.OrderStatus {
	switch status {
	case service.OrderStatusNew:
		return pb.OrderStatus_NEW
	case service.OrderStatusPending:
		return pb.OrderStatus_PENDING
	case service.OrderStatusProcessing:
		return pb.OrderStatus_PROCESSING
	case service.OrderStatusPaid:
		return pb.OrderStatus_PAID
	default:
		return pb.OrderStatus_UNKNOWN
	}
}
