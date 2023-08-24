package checkout

import (
	"context"

	"github.com/Mik3y-F/order-management-system/orders/internal/repository"
	"github.com/Mik3y-F/order-management-system/orders/internal/service"
	orders_pkg "github.com/Mik3y-F/order-management-system/orders/pkg"

	"github.com/Mik3y-F/order-management-system/payments/pkg/client"
	"github.com/Mik3y-F/order-management-system/pkg"
)

type CheckoutService struct {
	productRepository  repository.ProductRepository
	customerRepository repository.CustomerRepository
	orderRepository    repository.OrderRepository

	paymentsClient client.PaymentsClient
}

func NewCheckoutService(
	productRepository repository.ProductRepository,
	customerRepository repository.CustomerRepository,
	orderRepository repository.OrderRepository,
	paymentsClient client.PaymentsClient) *CheckoutService {

	return &CheckoutService{
		orderRepository:    orderRepository,
		productRepository:  productRepository,
		customerRepository: customerRepository,
		paymentsClient:     paymentsClient,
	}
}

func (s *CheckoutService) CheckPreconditions() {
	if s.productRepository == nil {
		panic("productRepository is required")
	}

	if s.customerRepository == nil {
		panic("customerRepository is required")
	}

	if s.orderRepository == nil {
		panic("orderRepository is required")
	}

	if s.paymentsClient == nil {
		panic("paymentsClient is required")
	}
}

func (s *CheckoutService) GetOrderCost(ctx context.Context, orderId string) (uint, error) {
	s.CheckPreconditions()

	order, err := s.orderRepository.GetOrder(ctx, orderId)
	if err != nil {
		return 0, err
	}

	var cost uint
	for _, item := range order.Items {
		product, err := s.productRepository.GetProduct(ctx, item.ProductId)
		if err != nil {
			return 0, service.Errorf(service.INTERNAL_ERROR, "failed to get product: %v", err)
		}

		cost += product.Price * item.Quantity
	}

	return cost, nil
}

func (s *CheckoutService) ProcessCheckout(ctx context.Context, orderId string) (*service.Order, error) {
	s.CheckPreconditions()

	order, err := s.orderRepository.UpdateOrderStatus(ctx, orderId, orders_pkg.OrderStatusProcessing)
	if err != nil {
		return nil, service.Errorf(service.INTERNAL_ERROR, "failed to update order status: %v", err)
	}

	customer, err := s.customerRepository.GetCustomer(ctx, order.CustomerId)
	if err != nil {
		return nil, service.Errorf(service.INTERNAL_ERROR, "failed to get customer: %v", err)
	}

	phoneNo, err := pkg.StringToUint(customer.Phone)
	if err != nil {
		return nil, service.Errorf(service.INTERNAL_ERROR, "failed to convert phone number to uint: %v", err)
	}

	cost, err := s.GetOrderCost(ctx, orderId)
	_, err = s.paymentsClient.ProcessMpesaPayment(ctx, &client.ProcessMpesaPaymentRequest{
		OrderId:     orderId,
		Amount:      uint32(cost),
		CustomerId:  order.CustomerId,
		PhoneNumber: uint64(phoneNo),
	})
	if err != nil {
		return nil, service.Errorf(service.INTERNAL_ERROR, "failed to process payment: %v", err)
	}

	order, err = s.orderRepository.GetOrder(ctx, orderId)
	if err != nil {
		return nil, service.Errorf(service.INTERNAL_ERROR, "failed to get order: %v", err)
	}

	return s.unmarshallRepositoryOrder(order), nil
}

func (s *CheckoutService) unmarshallOrderItem(item *repository.OrderItem) *service.OrderItem {
	return &service.OrderItem{
		Id:        item.Id,
		ProductId: item.ProductId,
		Quantity:  item.Quantity,
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
	}
}

func (s *CheckoutService) unmarshallRepositoryOrder(order *repository.Order) *service.Order {

	orderItems := make([]*service.OrderItem, len(order.Items))
	for i, item := range order.Items {
		orderItems[i] = s.unmarshallOrderItem(item)
	}

	return &service.Order{
		Id:          order.Id,
		CustomerId:  order.CustomerId,
		Items:       orderItems,
		OrderStatus: order.OrderStatus,
		CreatedAt:   order.CreatedAt,
		UpdatedAt:   order.UpdatedAt,
	}
}
