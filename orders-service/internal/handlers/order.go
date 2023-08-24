package handlers

import (
	"context"
	"fmt"

	pb "github.com/Mik3y-F/order-management-system/orders/api/generated"
	"github.com/Mik3y-F/order-management-system/orders/internal/repository"
	"github.com/Mik3y-F/order-management-system/orders/pkg"
)

func (s *GRPCServer) CreateOrder(ctx context.Context, in *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {

	var items []*repository.OrderItem
	for _, item := range in.GetOrderItems() {
		items = append(items, &repository.OrderItem{
			ProductId: item.GetProductId(),
			Quantity:  uint(item.GetQuantity()),
		})
	}

	p, err := s.OrderRepository.CreateOrder(ctx, &repository.Order{
		CustomerId: in.GetCustomerId(),
		Items:      items,
	})
	if err != nil {
		return nil, Error(fmt.Errorf("failed to create order: %w", err))
	}

	return &pb.CreateOrderResponse{
		Id: p.Id,
	}, nil
}

func (s *GRPCServer) GetOrder(ctx context.Context, in *pb.GetOrderRequest) (*pb.GetOrderResponse, error) {

	order, err := s.OrderRepository.GetOrder(ctx, in.GetId())
	if err != nil {
		return nil, Error(fmt.Errorf("failed to get order: %w", err))
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

	orders, err := s.OrderRepository.ListOrders(ctx)
	if err != nil {
		return nil, Error(fmt.Errorf("failed to list orders: %w", err))
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

func (s *GRPCServer) UpdateOrderStatus(
	ctx context.Context, in *pb.UpdateOrderStatusRequest) (*pb.UpdateOrderStatusResponse, error) {

	order, err := s.OrderRepository.UpdateOrderStatus(ctx, in.GetId(), getOrderStatus(in.GetStatus()))
	if err != nil {
		return nil, Error(fmt.Errorf("failed to update order status: %w", err))
	}

	return &pb.UpdateOrderStatusResponse{
		Id:         order.Id,
		CustomerId: order.CustomerId,
		Status:     getGRPCOrderStatus(order.OrderStatus),
	}, nil
}

func (s *GRPCServer) DeleteOrder(ctx context.Context, in *pb.DeleteOrderRequest) (*pb.DeleteOrderResponse, error) {

	err := s.OrderRepository.DeleteOrder(ctx, in.GetId())
	if err != nil {
		return nil, Error(fmt.Errorf("failed to delete order: %w", err))
	}

	return &pb.DeleteOrderResponse{}, nil
}

func (s *GRPCServer) CreateOrderItem(
	ctx context.Context, in *pb.CreateOrderItemRequest) (*pb.CreateOrderItemResponse, error) {

	orderItem, err := s.OrderRepository.CreateOrderItem(ctx, in.GetOrderId(), &repository.OrderItem{
		ProductId: in.GetProductId(),
		Quantity:  uint(in.GetQuantity()),
	})
	if err != nil {
		return nil, Error(fmt.Errorf("failed to create order item: %w", err))
	}

	return &pb.CreateOrderItemResponse{
		Id: orderItem.Id,
	}, nil
}

func (s *GRPCServer) GetOrderItem(ctx context.Context, in *pb.GetOrderItemRequest) (*pb.GetOrderItemResponse, error) {

	orderItem, err := s.OrderRepository.GetOrderItem(ctx, in.GetOrderId(), in.GetId())
	if err != nil {
		return nil, Error(fmt.Errorf("failed to get order item: %w", err))
	}

	return &pb.GetOrderItemResponse{
		Id:        orderItem.Id,
		ProductId: orderItem.ProductId,
		Quantity:  uint32(orderItem.Quantity),
	}, nil
}

func (s *GRPCServer) ListOrderItems(
	ctx context.Context, in *pb.ListOrderItemsRequest) (*pb.ListOrderItemsResponse, error) {

	orderItems, err := s.OrderRepository.ListOrderItems(ctx, in.GetOrderId())
	if err != nil {
		return nil, Error(fmt.Errorf("failed to list order items: %w", err))
	}

	var responseOrderItems []*pb.OrderItem
	for _, item := range orderItems {
		responseOrderItems = append(responseOrderItems, &pb.OrderItem{
			Id:        item.Id,
			ProductId: item.ProductId,
			Quantity:  uint32(item.Quantity),
		})
	}

	return &pb.ListOrderItemsResponse{
		OrderItems: responseOrderItems,
	}, nil
}

func (s *GRPCServer) UpdateOrderItem(
	ctx context.Context, in *pb.UpdateOrderItemRequest) (*pb.UpdateOrderItemResponse, error) {

	var quantity *uint
	if in.GetUpdate().GetQuantity() > 0 {
		quantity = new(uint)
		*quantity = uint(in.GetUpdate().GetQuantity())
	}
	orderItem, err := s.OrderRepository.UpdateOrderItem(ctx, in.GetOrderId(), in.GetId(), &repository.OrderItemUpdate{
		Quantity: quantity,
	})
	if err != nil {
		return nil, Error(fmt.Errorf("failed to update order item: %w", err))
	}

	return &pb.UpdateOrderItemResponse{
		Id:        orderItem.Id,
		ProductId: orderItem.ProductId,
		Quantity:  uint32(orderItem.Quantity),
	}, nil
}

func (s *GRPCServer) DeleteOrderItem(
	ctx context.Context, in *pb.DeleteOrderItemRequest) (*pb.DeleteOrderItemResponse, error) {

	err := s.OrderRepository.DeleteOrderItem(ctx, in.GetOrderId(), in.GetId())
	if err != nil {
		return nil, Error(fmt.Errorf("failed to delete order item: %w", err))
	}

	return &pb.DeleteOrderItemResponse{}, nil
}

func getGRPCOrderStatus(status pkg.OrderStatus) pb.OrderStatus {
	switch status {
	case pkg.OrderStatusNew:
		return pb.OrderStatus_NEW
	case pkg.OrderStatusPending:
		return pb.OrderStatus_PENDING
	case pkg.OrderStatusProcessing:
		return pb.OrderStatus_PROCESSING
	case pkg.OrderStatusPaid:
		return pb.OrderStatus_PAID
	default:
		return pb.OrderStatus_UNKNOWN
	}
}

func getOrderStatus(status pb.OrderStatus) pkg.OrderStatus {
	switch status {
	case pb.OrderStatus_NEW:
		return pkg.OrderStatusNew
	case pb.OrderStatus_PENDING:
		return pkg.OrderStatusPending
	case pb.OrderStatus_PROCESSING:
		return pkg.OrderStatusProcessing
	case pb.OrderStatus_PAID:
		return pkg.OrderStatusPaid
	default:
		return pkg.OrderStatusNew
	}
}
