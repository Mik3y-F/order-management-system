package firebase_test

import (
	"context"
	"reflect"
	"testing"
	"time"

	db "github.com/Mik3y-F/order-management-system/orders/internal/firebase"
	"github.com/Mik3y-F/order-management-system/orders/internal/service"
)

func deleteTestProduct(t *testing.T, ctx context.Context, productService service.ProductService, id string) {
	err := productService.DeleteProduct(ctx, id)
	if err != nil {
		t.Fatalf("failed to delete product: %v", err)
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
			name: "Check Preconditions Failed - nil DB",
			fields: fields{
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
				CreatedAt:   time.Now().Format(time.RFC3339),
				UpdatedAt:   time.Now().Format(time.RFC3339),
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
			// Clean up the created product
			defer deleteTestProduct(t, ctx, productService, got.Id)

			// Ignore the ID in the comparison since it's unpredictable
			got.Id = ""
			tt.want.Id = ""

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ProductService.CreateProduct() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProductService_GetProduct(t *testing.T) {

	ctx := context.Background()

	firebase := db.NewFirebaseService()
	firestoreClient, err := firebase.GetApp().Firestore(ctx)
	if err != nil {
		t.Fatalf("failed to create firestore client: %v", err)
	}
	defer firestoreClient.Close()

	firestoreService := db.NewFirestoreService(firestoreClient)
	productService := db.NewProductService(firestoreService)

	p, err := productService.CreateProduct(ctx, &service.Product{
		Name:        "Test Product",
		Description: "Test Description",
		Price:       100,
	})
	if err != nil {
		t.Fatalf("failed to create product: %v", err)
	}
	defer deleteTestProduct(t, ctx, productService, p.Id)

	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		args    args
		want    *service.Product
		wantErr bool
	}{
		{
			name: "Get Product Success",
			args: args{
				ctx: context.Background(),
				id:  p.Id,
			},
			want:    p,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := productService.GetProduct(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProductService.GetProduct() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ProductService.GetProduct() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProductService_ListProducts(t *testing.T) {

	ctx := context.Background()

	firebase := db.NewFirebaseService()
	firestoreClient, err := firebase.GetApp().Firestore(ctx)
	if err != nil {
		t.Fatalf("failed to create firestore client: %v", err)
	}
	defer firestoreClient.Close()

	firestoreService := db.NewFirestoreService(firestoreClient)
	productService := db.NewProductService(firestoreService)

	p, err := productService.CreateProduct(ctx, &service.Product{
		Name:        "Test Product",
		Description: "Test Description",
		Price:       100,
	})
	if err != nil {
		t.Fatalf("failed to create product: %v", err)
	}

	defer deleteTestProduct(t, ctx, productService, p.Id)

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    []*service.Product
		wantErr bool
	}{
		{
			name: "List Products Success",
			args: args{
				ctx: context.Background(),
			},
			want: []*service.Product{
				p,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := productService.ListProducts(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProductService.ListProducts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ProductService.ListProducts() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProductService_UpdateProduct(t *testing.T) {

	ctx := context.Background()

	firebase := db.NewFirebaseService()
	firestoreClient, err := firebase.GetApp().Firestore(ctx)
	if err != nil {
		t.Fatalf("failed to create firestore client: %v", err)
	}
	defer firestoreClient.Close()

	firestoreService := db.NewFirestoreService(firestoreClient)
	productService := db.NewProductService(firestoreService)

	p, err := productService.CreateProduct(ctx, &service.Product{
		Name:        "Test Product",
		Description: "Test Description",
		Price:       100,
	})
	if err != nil {
		t.Fatalf("failed to create product: %v", err)
	}

	defer deleteTestProduct(t, ctx, productService, p.Id)

	type args struct {
		ctx    context.Context
		id     string
		update *service.ProductUpdate
	}
	tests := []struct {
		name    string
		args    args
		want    *service.Product
		wantErr bool
	}{
		{
			name: "Update Product Success",
			args: args{
				ctx: context.Background(),
				id:  p.Id,
				update: &service.ProductUpdate{
					Name:        "Updated Test Product",
					Description: "Updated Test Description",
					Price:       200,
				},
			},
			want: &service.Product{
				Id:          p.Id,
				Name:        "Updated Test Product",
				Description: "Updated Test Description",
				Price:       200,
				CreatedAt:   p.CreatedAt,
				UpdatedAt:   time.Now().Format(time.RFC3339),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := productService.UpdateProduct(tt.args.ctx, tt.args.id, tt.args.update)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProductService.UpdateProduct() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ProductService.UpdateProduct() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProductService_DeleteProduct(t *testing.T) {

	ctx := context.Background()

	firebase := db.NewFirebaseService()
	firestoreClient, err := firebase.GetApp().Firestore(ctx)
	if err != nil {
		t.Fatalf("failed to create firestore client: %v", err)
	}
	defer firestoreClient.Close()

	firestoreService := db.NewFirestoreService(firestoreClient)
	productService := db.NewProductService(firestoreService)

	p, err := productService.CreateProduct(ctx, &service.Product{
		Name:        "Test Product",
		Description: "Test Description",
		Price:       100,
	})
	if err != nil {
		t.Fatalf("failed to create product: %v", err)
	}
	defer deleteTestProduct(t, ctx, productService, p.Id)

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
			name: "Delete Product Success",
			args: args{
				ctx: context.Background(),
				id:  p.Id,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := productService.DeleteProduct(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("ProductService.DeleteProduct() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
