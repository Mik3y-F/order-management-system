package mock

import (
	"context"

	"github.com/Mik3y-F/order-management-system/orders/internal/service"
)

var _ service.OrderService = (*OrderService)(nil)

type OrderService struct {
	CreateOrderFunc func(ctx context.Context, order *service.Order) (*service.Order, error)
	GetOrderFunc    func(ctx context.Context, id string) (*service.Order, error)
	ListOrdersFunc  func(ctx context.Context) ([]*service.Order, error)
	DeleteOrderFunc func(ctx context.Context, id string) error

	CreateOrderItemFunc func(ctx context.Context, orderId string, item *service.OrderItem) (*service.OrderItem, error)
	GetOrderItemFunc    func(ctx context.Context, orderId string, itemId string) (*service.OrderItem, error)
	ListOrderItemsFunc  func(ctx context.Context, orderId string) ([]*service.OrderItem, error)
	UpdateOrderItemFunc func(
		ctx context.Context, orderId string, itemId string, update *service.OrderItemUpdate) (*service.OrderItem, error)
	DeleteOrderItemFunc func(ctx context.Context, orderId string, itemId string) error
}

func (m *OrderService) CreateOrder(ctx context.Context, order *service.Order) (*service.Order, error) {
	return m.CreateOrderFunc(ctx, order)
}

func (m *OrderService) GetOrder(ctx context.Context, id string) (*service.Order, error) {
	return m.GetOrderFunc(ctx, id)
}

func (m *OrderService) ListOrders(ctx context.Context) ([]*service.Order, error) {
	return m.ListOrdersFunc(ctx)
}

func (m *OrderService) DeleteOrder(ctx context.Context, id string) error {
	return m.DeleteOrderFunc(ctx, id)
}

func (m *OrderService) CreateOrderItem(
	ctx context.Context, orderId string, item *service.OrderItem) (*service.OrderItem, error) {
	return m.CreateOrderItemFunc(ctx, orderId, item)
}

func (m *OrderService) GetOrderItem(ctx context.Context, orderId string, itemId string) (*service.OrderItem, error) {
	return m.GetOrderItemFunc(ctx, orderId, itemId)
}

func (m *OrderService) ListOrderItems(ctx context.Context, orderId string) ([]*service.OrderItem, error) {
	return m.ListOrderItemsFunc(ctx, orderId)
}

func (m *OrderService) UpdateOrderItem(
	ctx context.Context, orderId string, itemId string, update *service.OrderItemUpdate) (*service.OrderItem, error) {
	return m.UpdateOrderItemFunc(ctx, orderId, itemId, update)
}

func (m *OrderService) DeleteOrderItem(ctx context.Context, orderId string, itemId string) error {
	return m.DeleteOrderItemFunc(ctx, orderId, itemId)
}
