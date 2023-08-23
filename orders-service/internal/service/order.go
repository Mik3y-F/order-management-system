package service

import "context"

type OrderItem struct {
	Id        string `json:"id"`
	ProductId string `json:"product_id"`
	Quantity  uint   `json:"quantity"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type OrderItemUpdate struct {
	Quantity *uint `json:"quantity"`
}

type OrderStatus string

const (
	OrderStatusNew        OrderStatus = "new"
	OrderStatusPending    OrderStatus = "pending"
	OrderStatusProcessing OrderStatus = "processing"
	OrderStatusPaid       OrderStatus = "paid"
)

type Order struct {
	Id          string       `json:"id"`
	CustomerId  string       `json:"customer_id"`
	Items       []*OrderItem `json:"items"`
	OrderStatus OrderStatus  `json:"order_status"`
	CreatedAt   string       `json:"created_at"`
	UpdatedAt   string       `json:"updated_at"`
}

type OrderService interface {

	// Order CRUD
	CreateOrder(ctx context.Context, order *Order) (*Order, error)
	GetOrder(ctx context.Context, id string) (*Order, error)
	ListOrders(ctx context.Context) ([]*Order, error)
	UpdateOrderStatus(ctx context.Context, orderId string, status OrderStatus) (*Order, error)
	DeleteOrder(ctx context.Context, id string) error

	// OrderItem CRUD
	CreateOrderItem(ctx context.Context, orderId string, orderItem *OrderItem) (*OrderItem, error)
	GetOrderItem(ctx context.Context, orderId string, orderItemId string) (*OrderItem, error)
	ListOrderItems(ctx context.Context, orderId string) ([]*OrderItem, error)
	UpdateOrderItem(ctx context.Context, orderId string, orderItemId string, update *OrderItemUpdate) (*OrderItem, error)
	DeleteOrderItem(ctx context.Context, orderId string, orderItemId string) error
}
