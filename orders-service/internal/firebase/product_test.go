package firebase_test

import (
	"context"
	"reflect"
	"testing"

	db "github.com/Mik3y-F/order-management-system/orders/internal/firebase"
	"github.com/Mik3y-F/order-management-system/orders/internal/service"
)

func TestProductService_CreateProduct(t *testing.T) {

	ctx := context.Background()

	firebase := db.NewFirebaseService()
	firestoreClient, err := firebase.GetApp().Firestore(ctx)
	if err != nil {
		t.Fatalf("failed to create firestore client: %v", err)
	}
	defer firestoreClient.Close()

	firestoreService := db.NewFirestoreService(firestoreClient)
	productService := db.NewProductService(firestoreService)

	type args struct {
		ctx     context.Context
		product *service.Product
	}
	tests := []struct {
		name    string
		args    args
		want    *service.Product
		wantErr bool
	}{
		{
			name: "Create Product Success",
			args: args{
				ctx: context.Background(),
				product: &service.Product{
					Name:        "Test Product",
					Description: "Test Description",
					Price:       100,
				},
			},
			want: &service.Product{
				Name:        "Test Product",
				Description: "Test Description",
				Price:       100,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := productService.CreateProduct(tt.args.ctx, tt.args.product)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProductService.CreateProduct() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// Ignore the ID in the comparison since it's unpredictable
			got.Id = ""
			tt.want.Id = ""

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ProductService.CreateProduct() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProductService_CheckPreconditions(t *testing.T) {
	type fields struct {
		db *db.FirestoreService
	}
	tests := []struct {
		name      string
		fields    fields
		wantPanic bool
	}{
		{
			name:      "Check Preconditions Failed - nil DB",
			fields:    fields{
				db: nil,
			},
			wantPanic: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := db.NewProductService(tt.fields.db)
			defer func() {
				r := recover()
				if (r != nil) != tt.wantPanic {
					t.Errorf("ProductService.CheckPreconditions() panic = %v, wantPanic %v", r, tt.wantPanic)
				}
			}()
			s.CheckPreconditions()
		})
	}
}
