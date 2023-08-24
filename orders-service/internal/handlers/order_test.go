package handlers_test

import (
	"context"
	"reflect"
	"testing"

	pb "github.com/Mik3y-F/order-management-system/orders/api/generated"
	"github.com/Mik3y-F/order-management-system/orders/internal/service"
	"github.com/Mik3y-F/order-management-system/orders/pkg"
)

const (
	ERROR_ORDER_TRIGGER = "Error Order"
)

func mockCreateOrderFunc(ctx context.Context, o *service.Order) (*service.Order, error) {
	if o.CustomerId == ERROR_ORDER_TRIGGER {
		return nil, service.Errorf(service.INVALID_ERROR, "intentional error")
	}

	var items []*service.OrderItem
	for _, item := range o.Items {
		items = append(items, &service.OrderItem{
			ProductId: item.ProductId,
			Quantity:  uint(item.Quantity),
		})
	}

	return &service.Order{
		Id:         "1",
		CustomerId: o.CustomerId,
		Items:      items,
	}, nil
}
func TestGRPCServer_CreateOrder(t *testing.T) {
	s := NewTestGRPCServer(t)

	s.OrderService.CreateOrderFunc = mockCreateOrderFunc

	type args struct {
		ctx context.Context
		in  *pb.CreateOrderRequest
	}
	tests := []struct {
		name string

		args    args
		want    *pb.CreateOrderResponse
		wantErr bool
	}{
		{
			name: "Create Order Success",
			args: args{
				ctx: context.Background(),
				in: &pb.CreateOrderRequest{
					CustomerId: "1",
					OrderItems: []*pb.OrderItem{
						{
							ProductId: "1",
							Quantity:  1,
						},
					},
				},
			},
			want: &pb.CreateOrderResponse{
				Id: "1",
			},
		},
		{
			name: "Create Order Error",
			args: args{
				ctx: context.Background(),
				in: &pb.CreateOrderRequest{
					CustomerId: ERROR_ORDER_TRIGGER,
					OrderItems: []*pb.OrderItem{
						{
							ProductId: "1",
							Quantity:  1,
						},
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.CreateOrder(tt.args.ctx, tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("GRPCServer.CreateOrder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GRPCServer.CreateOrder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func mockGetOrderFunc(ctx context.Context, id string) (*service.Order, error) {
	if id == ERROR_ORDER_TRIGGER {
		return nil, service.Errorf(service.INVALID_ERROR, "intentional error")
	}

	return &service.Order{
		Id:          "1",
		CustomerId:  "1",
		OrderStatus: pkg.OrderStatusNew,
		Items: []*service.OrderItem{
			{
				ProductId: "1",
				Quantity:  1,
			},
		},
	}, nil
}
func TestGRPCServer_GetOrder(t *testing.T) {
	s := NewTestGRPCServer(t)

	s.OrderService.GetOrderFunc = mockGetOrderFunc

	type args struct {
		ctx context.Context
		in  *pb.GetOrderRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *pb.GetOrderResponse
		wantErr bool
	}{
		{
			name: "Get Order Success",
			args: args{
				ctx: context.Background(),
				in: &pb.GetOrderRequest{
					Id: "1",
				},
			},
			want: &pb.GetOrderResponse{
				Id:         "1",
				CustomerId: "1",
				Status:     pb.OrderStatus_NEW,
				OrderItems: []*pb.OrderItem{
					{
						ProductId: "1",
						Quantity:  1,
					},
				},
			},
		},
		{
			name: "Get Order Error",
			args: args{
				ctx: context.Background(),
				in: &pb.GetOrderRequest{
					Id: ERROR_ORDER_TRIGGER,
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := s.GetOrder(tt.args.ctx, tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("GRPCServer.GetOrder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GRPCServer.GetOrder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func mockListOrdersFunc(ctx context.Context) ([]*service.Order, error) {
	return []*service.Order{
		{
			Id:          "1",
			CustomerId:  "1",
			OrderStatus: pkg.OrderStatusNew,
			Items: []*service.OrderItem{
				{
					ProductId: "1",
					Quantity:  1,
				},
			},
		},
	}, nil
}

func TestGRPCServer_ListOrders(t *testing.T) {

	s := NewTestGRPCServer(t)

	s.OrderService.ListOrdersFunc = mockListOrdersFunc

	type args struct {
		ctx context.Context
		in  *pb.ListOrdersRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *pb.ListOrdersResponse
		wantErr bool
	}{
		{
			name: "List Orders Success",
			args: args{
				ctx: context.Background(),
				in:  &pb.ListOrdersRequest{},
			},
			want: &pb.ListOrdersResponse{
				Orders: []*pb.Order{
					{
						Id:         "1",
						CustomerId: "1",
						Status:     pb.OrderStatus_NEW,
						OrderItems: []*pb.OrderItem{
							{
								ProductId: "1",
								Quantity:  1,
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.ListOrders(tt.args.ctx, tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("GRPCServer.ListOrders() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GRPCServer.ListOrders() = %v, want %v", got, tt.want)
			}
		})
	}
}

func mockUpdateOrderStatusFunc(ctx context.Context, id string, status pkg.OrderStatus) (*service.Order, error) {
	if id == ERROR_ORDER_TRIGGER {
		return nil, service.Errorf(service.INVALID_ERROR, "intentional error")
	}

	return &service.Order{
		Id:          "1",
		CustomerId:  "1",
		OrderStatus: status,
		Items: []*service.OrderItem{
			{
				ProductId: "1",
				Quantity:  1,
			},
		},
	}, nil
}

func TestGRPCServer_UpdateOrderStatus(t *testing.T) {
	s := NewTestGRPCServer(t)

	s.OrderService.UpdateOrderStatusFunc = mockUpdateOrderStatusFunc

	type args struct {
		ctx context.Context
		in  *pb.UpdateOrderStatusRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *pb.UpdateOrderStatusResponse
		wantErr bool
	}{
		{
			name: "Update Order Status Success",
			args: args{
				ctx: context.Background(),
				in: &pb.UpdateOrderStatusRequest{
					Id:     "1",
					Status: pb.OrderStatus_NEW,
				},
			},
			want: &pb.UpdateOrderStatusResponse{
				Id:         "1",
				CustomerId: "1",
				Status:     pb.OrderStatus_NEW,
			},
		},
		{
			name: "Update Order Status Error",
			args: args{
				ctx: context.Background(),
				in: &pb.UpdateOrderStatusRequest{
					Id:     ERROR_ORDER_TRIGGER,
					Status: pb.OrderStatus_NEW,
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := s.UpdateOrderStatus(tt.args.ctx, tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("GRPCServer.UpdateOrderStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GRPCServer.UpdateOrderStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func mockDeleteOrderFunc(ctx context.Context, id string) error {
	if id == ERROR_ORDER_TRIGGER {
		return service.Errorf(service.INVALID_ERROR, "intentional error")
	}

	return nil
}

func TestGRPCServer_DeleteOrder(t *testing.T) {

	s := NewTestGRPCServer(t)

	s.OrderService.DeleteOrderFunc = mockDeleteOrderFunc

	type args struct {
		ctx context.Context
		in  *pb.DeleteOrderRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *pb.DeleteOrderResponse
		wantErr bool
	}{
		{
			name: "Delete Order Success",
			args: args{
				ctx: context.Background(),
				in: &pb.DeleteOrderRequest{
					Id: "1",
				},
			},
			want: &pb.DeleteOrderResponse{},
		},
		{
			name: "Delete Order Error",
			args: args{
				ctx: context.Background(),
				in: &pb.DeleteOrderRequest{
					Id: ERROR_ORDER_TRIGGER,
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := s.DeleteOrder(tt.args.ctx, tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("GRPCServer.DeleteOrder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GRPCServer.DeleteOrder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func mockCreateOrderItemFunc(ctx context.Context, orderId string, item *service.OrderItem) (*service.OrderItem, error) {
	if orderId == ERROR_ORDER_TRIGGER {
		return nil, service.Errorf(service.INVALID_ERROR, "intentional error")
	}

	return &service.OrderItem{
		Id:        "1",
		ProductId: item.ProductId,
		Quantity:  item.Quantity,
	}, nil
}

func TestGRPCServer_CreateOrderItem(t *testing.T) {

	s := NewTestGRPCServer(t)

	s.OrderService.CreateOrderItemFunc = mockCreateOrderItemFunc

	type args struct {
		ctx context.Context
		in  *pb.CreateOrderItemRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *pb.CreateOrderItemResponse
		wantErr bool
	}{
		{
			name: "Create Order Item Success",
			args: args{
				ctx: context.Background(),
				in: &pb.CreateOrderItemRequest{
					OrderId:   "1",
					ProductId: "1",
					Quantity:  1,
				},
			},
			want: &pb.CreateOrderItemResponse{
				Id: "1",
			},
			wantErr: false,
		},
		{
			name: "Create Order Item Error",
			args: args{
				ctx: context.Background(),
				in: &pb.CreateOrderItemRequest{
					OrderId:   ERROR_ORDER_TRIGGER,
					ProductId: "1",
					Quantity:  1,
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := s.CreateOrderItem(tt.args.ctx, tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("GRPCServer.CreateOrderItem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GRPCServer.CreateOrderItem() = %v, want %v", got, tt.want)
			}
		})
	}
}

func mockGetOrderItemFunc(ctx context.Context, orderId, itemId string) (*service.OrderItem, error) {
	if orderId == ERROR_ORDER_TRIGGER {
		return nil, service.Errorf(service.INVALID_ERROR, "intentional error")
	}

	return &service.OrderItem{
		Id:        "1",
		ProductId: "1",
		Quantity:  1,
	}, nil
}

func TestGRPCServer_GetOrderItem(t *testing.T) {
	s := NewTestGRPCServer(t)

	s.OrderService.GetOrderItemFunc = mockGetOrderItemFunc

	type args struct {
		ctx context.Context
		in  *pb.GetOrderItemRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *pb.GetOrderItemResponse
		wantErr bool
	}{
		{
			name: "Get Order Item Success",
			args: args{
				ctx: context.Background(),
				in: &pb.GetOrderItemRequest{
					OrderId: "1",
					Id:      "1",
				},
			},
			want: &pb.GetOrderItemResponse{
				Id:        "1",
				ProductId: "1",
				Quantity:  1,
			},
		},
		{
			name: "Get Order Item Error",
			args: args{
				ctx: context.Background(),
				in: &pb.GetOrderItemRequest{
					OrderId: ERROR_ORDER_TRIGGER,
					Id:      "1",
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := s.GetOrderItem(tt.args.ctx, tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("GRPCServer.GetOrderItem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GRPCServer.GetOrderItem() = %v, want %v", got, tt.want)
			}
		})
	}
}

func mockListOrderItemsFunc(ctx context.Context, orderId string) ([]*service.OrderItem, error) {
	if orderId == ERROR_ORDER_TRIGGER {
		return nil, service.Errorf(service.INVALID_ERROR, "intentional error")
	}

	return []*service.OrderItem{
		{
			Id:        "1",
			ProductId: "1",
			Quantity:  1,
		},
	}, nil
}

func TestGRPCServer_ListOrderItems(t *testing.T) {

	s := NewTestGRPCServer(t)
	s.OrderService.ListOrderItemsFunc = mockListOrderItemsFunc

	type args struct {
		ctx context.Context
		in  *pb.ListOrderItemsRequest
	}
	tests := []struct {
		name string

		args    args
		want    *pb.ListOrderItemsResponse
		wantErr bool
	}{
		{
			name: "List Order Items Success",
			args: args{
				ctx: context.Background(),
				in: &pb.ListOrderItemsRequest{
					OrderId: "1",
				},
			},
			want: &pb.ListOrderItemsResponse{
				OrderItems: []*pb.OrderItem{
					{
						Id:        "1",
						ProductId: "1",
						Quantity:  1,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "List Order Items Error",
			args: args{
				ctx: context.Background(),
				in: &pb.ListOrderItemsRequest{
					OrderId: ERROR_ORDER_TRIGGER,
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := s.ListOrderItems(tt.args.ctx, tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("GRPCServer.ListOrderItems() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GRPCServer.ListOrderItems() = %v, want %v", got, tt.want)
			}
		})
	}
}

func mockUpdateOrderItemFunc(ctx context.Context, orderId, itemId string, item *service.OrderItemUpdate) (*service.OrderItem, error) {
	if orderId == ERROR_ORDER_TRIGGER {
		return nil, service.Errorf(service.INVALID_ERROR, "intentional error")
	}

	return &service.OrderItem{
		Id:        "1",
		ProductId: "1",
		Quantity:  *item.Quantity,
	}, nil
}

func TestGRPCServer_UpdateOrderItem(t *testing.T) {

	s := NewTestGRPCServer(t)
	s.OrderService.UpdateOrderItemFunc = mockUpdateOrderItemFunc

	type args struct {
		ctx context.Context
		in  *pb.UpdateOrderItemRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *pb.UpdateOrderItemResponse
		wantErr bool
	}{
		{
			name: "Update Order Item Success",
			args: args{
				ctx: context.Background(),
				in: &pb.UpdateOrderItemRequest{
					OrderId: "1",
					Id:      "1",
					Update: &pb.OrderItemUpdate{
						Quantity: 6,
					},
				},
			},
			want: &pb.UpdateOrderItemResponse{
				Id:        "1",
				ProductId: "1",
				Quantity:  6,
			},
			wantErr: false,
		},
		{
			name: "Update Order Item Error",
			args: args{
				ctx: context.Background(),
				in: &pb.UpdateOrderItemRequest{
					OrderId: ERROR_ORDER_TRIGGER,
					Id:      "1",
					Update: &pb.OrderItemUpdate{
						Quantity: 6,
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.UpdateOrderItem(tt.args.ctx, tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("GRPCServer.UpdateOrderItem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GRPCServer.UpdateOrderItem() = %v, want %v", got, tt.want)
			}
		})
	}
}

func mockDeleteOrderItemFunc(ctx context.Context, orderId, itemId string) error {
	if orderId == ERROR_ORDER_TRIGGER {
		return service.Errorf(service.INVALID_ERROR, "intentional error")
	}

	return nil
}

func TestGRPCServer_DeleteOrderItem(t *testing.T) {

	s := NewTestGRPCServer(t)
	s.OrderService.DeleteOrderItemFunc = mockDeleteOrderItemFunc

	type args struct {
		ctx context.Context
		in  *pb.DeleteOrderItemRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *pb.DeleteOrderItemResponse
		wantErr bool
	}{
		{
			name: "Delete Order Item Success",
			args: args{
				ctx: context.Background(),
				in: &pb.DeleteOrderItemRequest{
					OrderId: "1",
					Id:      "1",
				},
			},
			want:    &pb.DeleteOrderItemResponse{},
			wantErr: false,
		},
		{
			name: "Delete Order Item Error",
			args: args{
				ctx: context.Background(),
				in: &pb.DeleteOrderItemRequest{
					OrderId: ERROR_ORDER_TRIGGER,
					Id:      "1",
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.DeleteOrderItem(tt.args.ctx, tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("GRPCServer.DeleteOrderItem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GRPCServer.DeleteOrderItem() = %v, want %v", got, tt.want)
			}
		})
	}
}
