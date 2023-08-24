package service

import (
	"context"

	"github.com/Mik3y-F/order-management-system/orders/pkg"
)

type OrderItem struct {
	Id        string `json:"id"`
	ProductId string `json:"product_id"`
	Quantity  uint   `json:"quantity"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func (o *OrderItem) Validate() error {
	if o.ProductId == "" {
		return Errorf(INVALID_ERROR, "product_id is required")
	}

	if o.Quantity == 0 {
		return Errorf(INVALID_ERROR, "quantity is required")
	}

	return nil
}

type OrderItemUpdate struct {
	Quantity *uint `json:"quantity"`
}

type Order struct {
	Id          string          `json:"id"`
	CustomerId  string          `json:"customer_id"`
	Items       []*OrderItem    `json:"items"`
	OrderStatus pkg.OrderStatus `json:"order_status"`
	CreatedAt   string          `json:"created_at"`
	UpdatedAt   string          `json:"updated_at"`
}

func (o *Order) Validate() error {
	if o.CustomerId == "" {
		return Errorf(INVALID_ERROR, "customer_id is required")
	}

	if len(o.Items) == 0 {
		return Errorf(INVALID_ERROR, "items are required")
	}

	for _, item := range o.Items {
		err := item.Validate()
		if err != nil {
			return err
		}
	}

	return nil
}

type OrderService interface {

	// Order CRUD
	CreateOrder(ctx context.Context, order *Order) (*Order, error)
	GetOrder(ctx context.Context, id string) (*Order, error)
	ListOrders(ctx context.Context) ([]*Order, error)
	UpdateOrderStatus(ctx context.Context, orderId string, status pkg.OrderStatus) (*Order, error)
	DeleteOrder(ctx context.Context, id string) error

	// OrderItem CRUD
	CreateOrderItem(ctx context.Context, orderId string, orderItem *OrderItem) (*OrderItem, error)
	GetOrderItem(ctx context.Context, orderId string, orderItemId string) (*OrderItem, error)
	ListOrderItems(ctx context.Context, orderId string) ([]*OrderItem, error)
	UpdateOrderItem(ctx context.Context, orderId string, orderItemId string, update *OrderItemUpdate) (*OrderItem, error)
	DeleteOrderItem(ctx context.Context, orderId string, orderItemId string) error
}
