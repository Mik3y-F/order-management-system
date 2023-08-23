package firebase_test

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"
	"time"

	db "github.com/Mik3y-F/order-management-system/orders/internal/firebase"
	"github.com/Mik3y-F/order-management-system/orders/internal/service"
)

func deleteTestOrder(t *testing.T, ctx context.Context, orderService service.OrderService, id string) {
	err := orderService.DeleteOrder(ctx, id)
	if err != nil {
		t.Fatalf("failed to delete order: %v", err)
	}
}

func TestOrderService_CheckPreconditions(t *testing.T) {
	type fields struct {
		db *db.FirestoreService
	}
	tests := []struct {
		name      string
		fields    fields
		wantPanic bool
	}{
		{
			name: "Check Preconditions Failed - nil DB",
			fields: fields{
				db: nil,
			},
			wantPanic: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := db.NewOrderService(tt.fields.db)
			defer func() {
				r := recover()
				if (r != nil) != tt.wantPanic {
					t.Errorf("OrderService.CheckPreconditions() panic = %v, wantPanic %v", r, tt.wantPanic)
				}
			}()
			s.CheckPreconditions()
		})
	}
}

func TestOrderService_CreateOrder(t *testing.T) {

	ctx := context.Background()

	firebase := db.NewFirebaseService()
	firestoreClient, err := firebase.GetApp().Firestore(ctx)
	if err != nil {
		t.Fatalf("failed to create firestore client: %v", err)
	}
	defer firestoreClient.Close()

	firestoreService := db.NewFirestoreService(firestoreClient)
	orderService := db.NewOrderService(firestoreService)

	type args struct {
		ctx   context.Context
		order *service.Order
	}
	tests := []struct {
		name    string
		args    args
		want    *service.Order
		wantErr bool
	}{
		{
			name: "Create Order Success",
			args: args{
				ctx: ctx,
				order: &service.Order{
					CustomerId: "customer-1",
					Items: []*service.OrderItem{
						{
							ProductId: "product-1",
							Quantity:  1,
						},
					},
				},
			},
			want: &service.Order{
				CustomerId: "customer-1",
				Items: []*service.OrderItem{
					{
						ProductId: "product-1",
						Quantity:  1,
						UpdatedAt: time.Now().Format(time.RFC3339),
						CreatedAt: time.Now().Format(time.RFC3339),
					},
				},
				CreatedAt: time.Now().Format(time.RFC3339),
				UpdatedAt: time.Now().Format(time.RFC3339),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := orderService.CreateOrder(tt.args.ctx, tt.args.order)

			if (err != nil) != tt.wantErr {
				t.Errorf("OrderService.CreateOrder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			defer deleteTestOrder(t, ctx, orderService, got.Id)

			// Clear out the fields that are set by the DB
			got.Id = ""
			tt.want.Id = ""

			for _, item := range got.Items {
				item.Id = ""
			}

			for _, item := range tt.want.Items {
				item.Id = ""
			}

			// Convert structs to JSON for easier comparison
			gotJSON, err := json.Marshal(got)
			if err != nil {
				t.Errorf("OrderService.CreateOrder() error = %v", err)
			}
			wantJSON, _ := json.Marshal(tt.want)
			if err != nil {
				t.Errorf("OrderService.CreateOrder() error = %v", err)
			}

			if string(gotJSON) != string(wantJSON) {
				t.Errorf("OrderService.CreateOrder() = %s, want %v", string(gotJSON), string(wantJSON))
			}
		})
	}
}

func TestOrderService_GetOrder(t *testing.T) {
	ctx := context.Background()

	firebase := db.NewFirebaseService()
	firestoreClient, err := firebase.GetApp().Firestore(ctx)
	if err != nil {
		t.Fatalf("failed to create firestore client: %v", err)
	}
	defer firestoreClient.Close()

	firestoreService := db.NewFirestoreService(firestoreClient)
	orderService := db.NewOrderService(firestoreService)

	testOrder, err := orderService.CreateOrder(ctx, &service.Order{
		CustomerId: "customer-1",
		Items: []*service.OrderItem{
			{
				ProductId: "product-1",
				Quantity:  1,
				UpdatedAt: time.Now().Format(time.RFC3339),
				CreatedAt: time.Now().Format(time.RFC3339),
			},
		},
	})
	if err != nil {
		t.Fatalf("failed to create test order: %v", err)
	}
	defer deleteTestOrder(t, ctx, orderService, testOrder.Id)

	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		args    args
		want    *service.Order
		wantErr bool
	}{
		{
			name: "Get Order Success",
			args: args{
				ctx: ctx,
				id:  testOrder.Id,
			},
			want:    testOrder,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := orderService.GetOrder(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("OrderService.GetOrder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// Convert structs to JSON for easier comparison
			gotJSON, err := json.Marshal(got)
			if err != nil {
				t.Errorf("OrderService.GetOrder() error = %v", err)
			}
			wantJSON, _ := json.Marshal(tt.want)
			if err != nil {
				t.Errorf("OrderService.GetOrder() error = %v", err)
			}

			if string(gotJSON) != string(wantJSON) {
				t.Errorf("OrderService.GetOrder() = %v, want %v", string(gotJSON), string(wantJSON))
			}
		})
	}
}

func TestOrderService_ListOrders(t *testing.T) {
	ctx := context.Background()

	firebase := db.NewFirebaseService()
	firestoreClient, err := firebase.GetApp().Firestore(ctx)
	if err != nil {
		t.Fatalf("failed to create firestore client: %v", err)
	}
	defer firestoreClient.Close()

	firestoreService := db.NewFirestoreService(firestoreClient)
	orderService := db.NewOrderService(firestoreService)

	testOrder, err := orderService.CreateOrder(ctx, &service.Order{
		CustomerId: "customer-1",
		Items: []*service.OrderItem{
			{
				ProductId: "product-1",
				Quantity:  1,
				UpdatedAt: time.Now().Format(time.RFC3339),
				CreatedAt: time.Now().Format(time.RFC3339),
			},
		},
	})
	if err != nil {
		t.Fatalf("failed to create test order: %v", err)
	}
	defer deleteTestOrder(t, ctx, orderService, testOrder.Id)

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    []*service.Order
		wantErr bool
	}{
		{
			name: "List Orders Success",
			args: args{
				ctx: ctx,
			},
			want: []*service.Order{
				testOrder,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := orderService.ListOrders(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("OrderService.ListOrders() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Convert structs to JSON for easier comparison
			gotJSON, err := json.Marshal(got)
			if err != nil {
				t.Errorf("OrderService.ListOrders() error = %v", err)
			}
			wantJSON, _ := json.Marshal(tt.want)
			if err != nil {
				t.Errorf("OrderService.ListOrders() error = %v", err)
			}

			if string(gotJSON) != string(wantJSON) {
				t.Errorf("OrderService.ListOrders() = %v, want %v", string(gotJSON), string(wantJSON))
			}
		})
	}
}

func TestOrderService_DeleteOrder(t *testing.T) {
	ctx := context.Background()

	firebase := db.NewFirebaseService()
	firestoreClient, err := firebase.GetApp().Firestore(ctx)
	if err != nil {
		t.Fatalf("failed to create firestore client: %v", err)
	}
	defer firestoreClient.Close()

	firestoreService := db.NewFirestoreService(firestoreClient)
	orderService := db.NewOrderService(firestoreService)

	testOrder, err := orderService.CreateOrder(ctx, &service.Order{
		CustomerId: "customer-1",
		Items: []*service.OrderItem{
			{
				ProductId: "product-1",
				Quantity:  1,
				UpdatedAt: time.Now().Format(time.RFC3339),
				CreatedAt: time.Now().Format(time.RFC3339),
			},
		},
	})
	if err != nil {
		t.Fatalf("failed to create test order: %v", err)
	}
	defer deleteTestOrder(t, ctx, orderService, testOrder.Id)

	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Delete Order Success",
			args:    args{ctx: ctx, id: testOrder.Id},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := orderService.DeleteOrder(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("OrderService.DeleteOrder() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOrderService_CreateOrderItem(t *testing.T) {

	ctx := context.Background()

	firebase := db.NewFirebaseService()
	firestoreClient, err := firebase.GetApp().Firestore(ctx)
	if err != nil {
		t.Fatalf("failed to create firestore client: %v", err)
	}
	defer firestoreClient.Close()

	firestoreService := db.NewFirestoreService(firestoreClient)
	orderService := db.NewOrderService(firestoreService)

	testOrder, err := orderService.CreateOrder(ctx, &service.Order{
		CustomerId: "customer-1",
		Items: []*service.OrderItem{
			{
				ProductId: "product-1",
				Quantity:  1,
				UpdatedAt: time.Now().Format(time.RFC3339),
				CreatedAt: time.Now().Format(time.RFC3339),
			},
		},
	})
	if err != nil {
		t.Fatalf("failed to create test order: %v", err)
	}
	defer deleteTestOrder(t, ctx, orderService, testOrder.Id)

	type args struct {
		ctx       context.Context
		orderId   string
		orderItem *service.OrderItem
	}
	tests := []struct {
		name    string
		args    args
		want    *service.OrderItem
		wantErr bool
	}{
		{
			name: "Create Order Item Success",
			args: args{
				ctx:     ctx,
				orderId: testOrder.Id,
				orderItem: &service.OrderItem{
					ProductId: "product-2",
					Quantity:  1,
				},
			},
			want: &service.OrderItem{
				ProductId: "product-2",
				Quantity:  1,
				UpdatedAt: time.Now().Format(time.RFC3339),
				CreatedAt: time.Now().Format(time.RFC3339),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// no need to delete the order item as it will be deleted when the order is deleted
			got, err := orderService.CreateOrderItem(tt.args.ctx, tt.args.orderId, tt.args.orderItem)
			if (err != nil) != tt.wantErr {
				t.Errorf("OrderService.CreateOrderItem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Clear out the fields that are set by the DB
			got.Id = ""
			tt.want.Id = ""

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OrderService.CreateOrderItem() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrderService_CreateOrderItems(t *testing.T) {
	ctx := context.Background()

	firebase := db.NewFirebaseService()
	firestoreClient, err := firebase.GetApp().Firestore(ctx)
	if err != nil {
		t.Fatalf("failed to create firestore client: %v", err)
	}
	defer firestoreClient.Close()

	firestoreService := db.NewFirestoreService(firestoreClient)
	orderService := db.NewOrderService(firestoreService)

	testOrder, err := orderService.CreateOrder(ctx, &service.Order{
		CustomerId: "customer-1",
		Items: []*service.OrderItem{
			{
				ProductId: "product-1",
				Quantity:  1,
				UpdatedAt: time.Now().Format(time.RFC3339),
				CreatedAt: time.Now().Format(time.RFC3339),
			},
		},
	})
	if err != nil {
		t.Fatalf("failed to create test order: %v", err)
	}
	defer deleteTestOrder(t, ctx, orderService, testOrder.Id)

	type args struct {
		ctx        context.Context
		orderId    string
		orderItems []*service.OrderItem
	}
	tests := []struct {
		name    string
		args    args
		want    []*service.OrderItem
		wantErr bool
	}{
		{
			name: "Create Order Items Success",
			args: args{
				ctx:     ctx,
				orderId: testOrder.Id,
				orderItems: []*service.OrderItem{
					{
						ProductId: "product-2",
						Quantity:  1,
					},
					{
						ProductId: "product-3",
						Quantity:  1,
					},
				},
			},
			want: []*service.OrderItem{
				{
					ProductId: "product-2",
					Quantity:  1,
					UpdatedAt: time.Now().Format(time.RFC3339),
					CreatedAt: time.Now().Format(time.RFC3339),
				},
				{
					ProductId: "product-3",
					Quantity:  1,
					UpdatedAt: time.Now().Format(time.RFC3339),
					CreatedAt: time.Now().Format(time.RFC3339),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := orderService.CreateOrderItems(tt.args.ctx, tt.args.orderId, tt.args.orderItems)
			if (err != nil) != tt.wantErr {
				t.Errorf("OrderService.CreateOrderItems() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			for _, item := range got {
				item.Id = ""
			}

			for _, item := range tt.want {
				item.Id = ""
			}

			// Convert structs to JSON for easier comparison
			gotJSON, err := json.Marshal(got)
			if err != nil {
				t.Errorf("OrderService.CreateOrderItems() error = %v", err)
			}
			wantJSON, _ := json.Marshal(tt.want)
			if err != nil {
				t.Errorf("OrderService.CreateOrderItems() error = %v", err)
			}

			if string(gotJSON) != string(wantJSON) {
				t.Errorf("OrderService.CreateOrderItems() = %v, want %v", string(gotJSON), string(wantJSON))
			}
		})
	}
}

func TestOrderService_GetOrderItem(t *testing.T) {
	ctx := context.Background()

	firebase := db.NewFirebaseService()
	firestoreClient, err := firebase.GetApp().Firestore(ctx)
	if err != nil {
		t.Fatalf("failed to create firestore client: %v", err)
	}
	defer firestoreClient.Close()

	firestoreService := db.NewFirestoreService(firestoreClient)
	orderService := db.NewOrderService(firestoreService)

	testOrder, err := orderService.CreateOrder(ctx, &service.Order{
		CustomerId: "customer-1",
		Items: []*service.OrderItem{
			{
				ProductId: "product-1",
				Quantity:  1,
				UpdatedAt: time.Now().Format(time.RFC3339),
				CreatedAt: time.Now().Format(time.RFC3339),
			},
		},
	})
	if err != nil {
		t.Fatalf("failed to create test order: %v", err)
	}
	defer deleteTestOrder(t, ctx, orderService, testOrder.Id)

	type args struct {
		ctx         context.Context
		orderId     string
		orderItemId string
	}
	tests := []struct {
		name string

		args    args
		want    *service.OrderItem
		wantErr bool
	}{
		{
			name: "Get Order Item Success",
			args: args{
				ctx:         ctx,
				orderId:     testOrder.Id,
				orderItemId: testOrder.Items[0].Id,
			},
			want: testOrder.Items[0],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := orderService.GetOrderItem(tt.args.ctx, tt.args.orderId, tt.args.orderItemId)
			if (err != nil) != tt.wantErr {
				t.Errorf("OrderService.GetOrderItem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OrderService.GetOrderItem() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrderService_ListOrderItems(t *testing.T) {
	ctx := context.Background()

	firebase := db.NewFirebaseService()
	firestoreClient, err := firebase.GetApp().Firestore(ctx)
	if err != nil {
		t.Fatalf("failed to create firestore client: %v", err)
	}
	defer firestoreClient.Close()

	firestoreService := db.NewFirestoreService(firestoreClient)
	orderService := db.NewOrderService(firestoreService)

	testOrder, err := orderService.CreateOrder(ctx, &service.Order{
		CustomerId: "customer-1",
		Items: []*service.OrderItem{
			{
				ProductId: "product-1",
				Quantity:  1,
				UpdatedAt: time.Now().Format(time.RFC3339),
				CreatedAt: time.Now().Format(time.RFC3339),
			},
		},
	})
	if err != nil {
		t.Fatalf("failed to create test order: %v", err)
	}
	defer deleteTestOrder(t, ctx, orderService, testOrder.Id)

	type args struct {
		ctx     context.Context
		orderId string
	}
	tests := []struct {
		name    string
		args    args
		want    []*service.OrderItem
		wantErr bool
	}{
		{
			name: "List Order Items Success",
			args: args{
				ctx:     ctx,
				orderId: testOrder.Id,
			},
			want: testOrder.Items,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := orderService.ListOrderItems(tt.args.ctx, tt.args.orderId)
			if (err != nil) != tt.wantErr {
				t.Errorf("OrderService.ListOrderItems() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OrderService.ListOrderItems() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrderService_UpdateOrderItem(t *testing.T) {
	ctx := context.Background()

	firebase := db.NewFirebaseService()
	firestoreClient, err := firebase.GetApp().Firestore(ctx)
	if err != nil {
		t.Fatalf("failed to create firestore client: %v", err)
	}
	defer firestoreClient.Close()

	firestoreService := db.NewFirestoreService(firestoreClient)
	orderService := db.NewOrderService(firestoreService)

	testOrder, err := orderService.CreateOrder(ctx, &service.Order{
		CustomerId: "customer-1",
		Items: []*service.OrderItem{
			{
				ProductId: "product-1",
				Quantity:  1,
				UpdatedAt: time.Now().Format(time.RFC3339),
				CreatedAt: time.Now().Format(time.RFC3339),
			},
		},
	})
	if err != nil {
		t.Fatalf("failed to create test order: %v", err)
	}
	defer deleteTestOrder(t, ctx, orderService, testOrder.Id)

	type args struct {
		ctx         context.Context
		orderId     string
		orderItemId string
		update      *service.OrderItemUpdate
	}
	tests := []struct {
		name    string
		args    args
		want    *service.OrderItem
		wantErr bool
	}{
		{
			name: "Update Order Item Success",
			args: args{
				ctx:         ctx,
				orderId:     testOrder.Id,
				orderItemId: testOrder.Items[0].Id,
				update: &service.OrderItemUpdate{
					Quantity: func(i uint) *uint { return &i }(2),
				},
			},
			want: &service.OrderItem{
				Id:        testOrder.Items[0].Id,
				ProductId: "product-1",
				Quantity:  2,
				UpdatedAt: time.Now().Format(time.RFC3339),
				CreatedAt: time.Now().Format(time.RFC3339),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := orderService.UpdateOrderItem(tt.args.ctx, tt.args.orderId, tt.args.orderItemId, tt.args.update)
			if (err != nil) != tt.wantErr {
				t.Errorf("OrderService.UpdateOrderItem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Convert structs to JSON for easier comparison
			gotJSON, err := json.Marshal(got)
			if err != nil {
				t.Errorf("OrderService.UpdateOrderItem() error = %v", err)
			}
			wantJSON, _ := json.Marshal(tt.want)
			if err != nil {
				t.Errorf("OrderService.UpdateOrderItem() error = %v", err)
			}

			if string(gotJSON) != string(wantJSON) {
				t.Errorf("OrderService.UpdateOrderItem() = %v, want %v", string(gotJSON), string(wantJSON))
			}
		})
	}
}

func TestOrderService_DeleteOrderItem(t *testing.T) {
	ctx := context.Background()

	firebase := db.NewFirebaseService()
	firestoreClient, err := firebase.GetApp().Firestore(ctx)
	if err != nil {
		t.Fatalf("failed to create firestore client: %v", err)
	}
	defer firestoreClient.Close()

	firestoreService := db.NewFirestoreService(firestoreClient)
	orderService := db.NewOrderService(firestoreService)

	testOrder, err := orderService.CreateOrder(ctx, &service.Order{
		CustomerId: "customer-1",
		Items: []*service.OrderItem{
			{
				ProductId: "product-1",
				Quantity:  1,
				UpdatedAt: time.Now().Format(time.RFC3339),
				CreatedAt: time.Now().Format(time.RFC3339),
			},
		},
	})
	if err != nil {
		t.Fatalf("failed to create test order: %v", err)
	}
	defer deleteTestOrder(t, ctx, orderService, testOrder.Id)

	type args struct {
		ctx         context.Context
		orderId     string
		orderItemId string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Delete Order Item Success",
			args: args{
				ctx:         ctx,
				orderId:     testOrder.Id,
				orderItemId: testOrder.Items[0].Id,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := orderService.DeleteOrderItem(tt.args.ctx, tt.args.orderId, tt.args.orderItemId); (err != nil) != tt.wantErr {
				t.Errorf("OrderService.DeleteOrderItem() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
