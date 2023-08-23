package firebase

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/Mik3y-F/order-management-system/orders/internal/service"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ service.OrderService = (*OrderService)(nil)

type OrderService struct {
	db *FirestoreService
}

func NewOrderService(db *FirestoreService) *OrderService {
	return &OrderService{
		db: db,
	}
}

func (s *OrderService) CheckPreconditions() {
	if s.db == nil {
		panic("no DB service provided")
	}
}

func (s *OrderService) orderCollection() *firestore.CollectionRef {
	s.CheckPreconditions()

	return s.db.client.Collection("orders")
}

func (s *OrderService) CreateOrder(ctx context.Context, order *service.Order) (*service.Order, error) {
	s.CheckPreconditions()

	// Set CreatedAt and UpdatedAt to the current time
	currentTime := time.Now()

	order.CreatedAt = currentTime.Format(time.RFC3339)
	order.UpdatedAt = currentTime.Format(time.RFC3339)

	err := order.Validate()
	if err != nil {
		return nil, service.Errorf(service.INVALID_ERROR, "invalid order details provided: %v", err)
	}

	orderModel := s.marshallOrder(order)

	docRef, _, err := s.orderCollection().Add(ctx, orderModel)
	if err != nil {
		return nil, service.Errorf(service.INTERNAL_ERROR, "failed to create order: %v", err)
	}

	order.Id = docRef.ID

	_, err = s.CreateOrderItems(ctx, order.Id, order.Items)
	if err != nil {
		return nil, service.Errorf(service.INTERNAL_ERROR, "failed to create order items: %v", err)
	}

	return order, nil
}

func (s *OrderService) GetOrder(ctx context.Context, id string) (*service.Order, error) {
	s.CheckPreconditions()

	if id == "" {
		return nil, service.Errorf(service.INVALID_ERROR, "id is required")
	}

	docRef, err := s.orderCollection().Doc(id).Get(ctx)
	if status.Code(err) == codes.NotFound {
		return nil, service.Errorf(service.NOT_FOUND_ERROR, "order not found")
	} else if err != nil {
		return nil, service.Errorf(service.INTERNAL_ERROR, "failed to get order: %v", err)
	}

	orderModel := &OrderModel{}
	if err := docRef.DataTo(orderModel); err != nil {
		return nil, service.Errorf(service.INTERNAL_ERROR, "failed to unmarshall order: %v", err)
	}

	order := s.unmarshallOrder(orderModel)

	order.Id = id

	orderItems, err := s.ListOrderItems(ctx, id)
	if err != nil {
		return nil, service.Errorf(service.INTERNAL_ERROR, "failed to get order items: %v", err)
	}

	order.Items = orderItems

	return order, nil
}

func (s *OrderService) ListOrders(ctx context.Context) ([]*service.Order, error) {
	s.CheckPreconditions()

	iter := s.orderCollection().Documents(ctx)

	orders := make([]*service.Order, 0)

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, service.Errorf(service.INTERNAL_ERROR, "failed to iterate orders: %v", err)
		}

		orderModel := &OrderModel{}
		if err := doc.DataTo(orderModel); err != nil {
			return nil, service.Errorf(service.INTERNAL_ERROR, "failed to unmarshall order: %v", err)
		}

		order := s.unmarshallOrder(orderModel)

		order.Id = doc.Ref.ID

		orderItems, err := s.ListOrderItems(ctx, order.Id)
		if err != nil {
			return nil, service.Errorf(service.INTERNAL_ERROR, "failed to get order items: %v", err)
		}

		order.Items = orderItems

		orders = append(orders, order)
	}

	return orders, nil
}

func (s *OrderService) UpdateOrderStatus(
	ctx context.Context, orderId string, status service.OrderStatus) (*service.Order, error) {
	s.CheckPreconditions()

	order, err := s.GetOrder(ctx, orderId)
	if err != nil {
		return nil, err
	}

	order.OrderStatus = status
	order.UpdatedAt = time.Now().Format(time.RFC3339)

	orderModel := s.marshallOrder(order)
	_, err = s.orderCollection().Doc(orderId).Set(ctx, orderModel)
	if err != nil {
		return nil, service.Errorf(service.INTERNAL_ERROR, "failed to update order status: %v", err)
	}

	return s.GetOrder(ctx, orderId)
}

func (s *OrderService) DeleteOrder(ctx context.Context, id string) error {
	s.CheckPreconditions()

	_, err := s.orderCollection().Doc(id).Delete(ctx)
	if err != nil {
		return service.Errorf(service.INTERNAL_ERROR, "failed to delete order: %v", err)
	}

	return nil
}

func (s *OrderService) orderItemCollection(orderId string) *firestore.CollectionRef {
	s.CheckPreconditions()

	return s.orderCollection().Doc(orderId).Collection("items")
}

func (s *OrderService) CreateOrderItem(
	ctx context.Context, orderId string, orderItem *service.OrderItem) (*service.OrderItem, error) {

	s.CheckPreconditions()

	// Set CreatedAt and UpdatedAt to the current time
	currentTime := time.Now()

	orderItem.CreatedAt = currentTime.Format(time.RFC3339)
	orderItem.UpdatedAt = currentTime.Format(time.RFC3339)

	err := orderItem.Validate()
	if err != nil {
		return nil, service.Errorf(service.INVALID_ERROR, "invalid order item details provided: %v", err)
	}

	orderItemModel := s.marshallOrderItem(orderItem)

	docRef, _, err := s.orderItemCollection(orderId).Add(ctx, orderItemModel)
	if err != nil {
		return nil, service.Errorf(service.INTERNAL_ERROR, "failed to create order item: %v", err)
	}

	orderItem.Id = docRef.ID

	return orderItem, nil
}

func (s *OrderService) CreateOrderItems(
	ctx context.Context, orderId string, orderItems []*service.OrderItem) ([]*service.OrderItem, error) {

	s.CheckPreconditions()

	if orderId == "" {
		return nil, service.Errorf(service.INVALID_ERROR, "order id is required")
	}

	bulkWriter := s.db.client.BulkWriter(ctx)

	currentTime := time.Now().Format(time.RFC3339)
	var createdOrderItems []*service.OrderItem

	for _, orderItem := range orderItems {
		// Set CreatedAt and UpdatedAt to the current time
		orderItem.CreatedAt = currentTime
		orderItem.UpdatedAt = currentTime

		orderItemModel := s.marshallOrderItem(orderItem)
		docRef := s.orderItemCollection(orderId).NewDoc() // Create a new document reference.

		orderItem.Id = docRef.ID
		createdOrderItems = append(createdOrderItems, orderItem)

		_, err := bulkWriter.Create(docRef, orderItemModel)
		if err != nil {
			return nil, service.Errorf(service.INTERNAL_ERROR, "failed to create order item: %v", err)
		}
	}

	bulkWriter.Flush()

	return createdOrderItems, nil
}

func (s *OrderService) GetOrderItem(
	ctx context.Context, orderId string, orderItemId string) (*service.OrderItem, error) {

	s.CheckPreconditions()

	if orderId == "" {
		return nil, service.Errorf(service.INVALID_ERROR, "order id is required")
	} else if orderItemId == "" {
		return nil, service.Errorf(service.INVALID_ERROR, "order item id is required")
	}

	docRef, err := s.orderItemCollection(orderId).Doc(orderItemId).Get(ctx)
	if status.Code(err) == codes.NotFound {
		return nil, service.Errorf(service.NOT_FOUND_ERROR, "order item not found")
	} else if err != nil {
		return nil, service.Errorf(service.INTERNAL_ERROR, "failed to get order item: %v", err)
	}

	orderItemModel := &OrderItemModel{}
	if err := docRef.DataTo(orderItemModel); err != nil {
		return nil, service.Errorf(service.INTERNAL_ERROR, "failed to unmarshall order item: %v", err)
	}

	orderItem := s.unmarshallOrderItem(orderItemModel)

	orderItem.Id = docRef.Ref.ID

	return orderItem, nil
}

func (s *OrderService) ListOrderItems(ctx context.Context, orderId string) ([]*service.OrderItem, error) {
	s.CheckPreconditions()

	if orderId == "" {
		return nil, service.Errorf(service.INVALID_ERROR, "order id is required")
	}

	iter := s.orderItemCollection(orderId).Documents(ctx)

	orderItems := make([]*service.OrderItem, 0)

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		orderItemModel := &OrderItemModel{}
		if err := doc.DataTo(orderItemModel); err != nil {
			return nil, service.Errorf(service.INTERNAL_ERROR, "failed to unmarshall order item: %v", err)
		}

		orderItem := s.unmarshallOrderItem(orderItemModel)

		orderItem.Id = doc.Ref.ID

		orderItems = append(orderItems, orderItem)
	}

	return orderItems, nil
}

func (s *OrderService) UpdateOrderItem(
	ctx context.Context, orderId string, orderItemId string, update *service.OrderItemUpdate) (*service.OrderItem, error) {

	s.CheckPreconditions()

	if orderId == "" {
		return nil, service.Errorf(service.INVALID_ERROR, "order id is required")
	} else if orderItemId == "" {
		return nil, service.Errorf(service.INVALID_ERROR, "order item id is required")
	}

	orderItem, err := s.GetOrderItem(ctx, orderId, orderItemId)
	if err != nil {
		return nil, err
	}

	if v := update.Quantity; v != nil {
		orderItem.Quantity = *v
	}

	// Set UpdatedAt to the current time
	orderItem.UpdatedAt = time.Now().Format(time.RFC3339)

	orderItemModel := s.marshallOrderItem(orderItem)
	_, err = s.orderItemCollection(orderId).Doc(orderItemId).Set(ctx, orderItemModel)
	if err != nil {
		return nil, service.Errorf(service.INTERNAL_ERROR, "failed to update order item: %v", err)
	}

	return s.GetOrderItem(ctx, orderId, orderItemId)
}

func (s *OrderService) DeleteOrderItem(ctx context.Context, orderId string, orderItemId string) error {
	s.CheckPreconditions()

	_, err := s.orderItemCollection(orderId).Doc(orderItemId).Delete(ctx)
	if err != nil {
		return service.Errorf(service.INTERNAL_ERROR, "failed to delete order item: %v", err)
	}

	return nil
}

func (s *OrderService) marshallOrder(order *service.Order) *OrderModel {
	return &OrderModel{
		CustomerId: order.CustomerId,
		Items:      s.marshallOrderItems(order.Items),
		CreatedAt:  order.CreatedAt,
		UpdatedAt:  order.UpdatedAt,
	}
}

func (s *OrderService) unmarshallOrder(order *OrderModel) *service.Order {
	return &service.Order{
		CustomerId: order.CustomerId,
		Items:      s.unmarshallOrderItems(order.Items),
		CreatedAt:  order.CreatedAt,
		UpdatedAt:  order.UpdatedAt,
	}
}

func (s *OrderService) marshallOrderItems(items []*service.OrderItem) []*OrderItemModel {
	orderItems := make([]*OrderItemModel, 0)

	for _, item := range items {
		orderItems = append(orderItems, s.marshallOrderItem(item))
	}

	return orderItems
}

func (s *OrderService) unmarshallOrderItems(items []*OrderItemModel) []*service.OrderItem {
	orderItems := make([]*service.OrderItem, 0)

	for _, item := range items {
		orderItems = append(orderItems, s.unmarshallOrderItem(item))
	}

	return orderItems
}

func (s *OrderService) marshallOrderItem(item *service.OrderItem) *OrderItemModel {
	return &OrderItemModel{
		ProductId: item.ProductId,
		Quantity:  int(item.Quantity),
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
	}
}

func (s *OrderService) unmarshallOrderItem(item *OrderItemModel) *service.OrderItem {
	return &service.OrderItem{
		ProductId: item.ProductId,
		Quantity:  uint(item.Quantity),
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
	}
}
