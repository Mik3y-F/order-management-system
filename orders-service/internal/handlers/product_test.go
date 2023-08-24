package handlers_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	pb "github.com/Mik3y-F/order-management-system/orders/api/generated"
	"github.com/Mik3y-F/order-management-system/orders/internal/repository"
)

const (
	ERROR_PRODUCT_TRIGGER = "Error Product"
)

func mockCreateProductFunc(ctx context.Context, p *repository.Product) (*repository.Product, error) {
	if p.Name == ERROR_PRODUCT_TRIGGER {
		return nil, errors.New("intentional error")
	}
	return &repository.Product{
		Id:          "1",
		Name:        "Test Product",
		Description: "Test Description",
		Price:       100,
	}, nil
}

func TestGRPCServer_CreateProduct(t *testing.T) {

	s := NewTestGRPCServer(t)

	s.ProductRepository.CreateProductFunc = mockCreateProductFunc

	tests := []struct {
		name    string
		request *pb.CreateProductRequest
		want    *pb.CreateProductResponse
		wantErr bool
	}{
		{
			name: "Create Product Success",
			request: &pb.CreateProductRequest{
				Name:        "Test Product",
				Description: "Test Description",
				Price:       100,
			},
			want: &pb.CreateProductResponse{
				Id: "1",
			},
			wantErr: false,
		},
		{
			name: "Create Product Error",
			request: &pb.CreateProductRequest{
				Name:        "Error Product",
				Description: "Test Description",
				Price:       100,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := s.GRPCServer.CreateProduct(context.Background(), tt.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("GRPCServer.CreateProduct() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GRPCServer.CreateProduct() = %v, want %v", got, tt.want)
			}
		})
	}
}

func mockGetProductFunc(ctx context.Context, id string) (*repository.Product, error) {
	if id == ERROR_PRODUCT_TRIGGER {
		return nil, errors.New("intentional error")
	}
	return &repository.Product{
		Id:          "1",
		Name:        "Test Product",
		Description: "Test Description",
		Price:       100,
	}, nil
}

func TestGRPCServer_GetProduct(t *testing.T) {

	s := NewTestGRPCServer(t)

	s.ProductRepository.GetProductFunc = mockGetProductFunc

	type args struct {
		ctx context.Context
		in  *pb.GetProductRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *pb.GetProductResponse
		wantErr bool
	}{
		{
			name: "Get Product Success",
			args: args{
				ctx: context.Background(),
				in: &pb.GetProductRequest{
					Id: "1",
				},
			},
			want: &pb.GetProductResponse{
				Id:          "1",
				Name:        "Test Product",
				Description: "Test Description",
				Price:       100,
			},
			wantErr: false,
		},
		{
			name: "Get Product Error",
			args: args{
				ctx: context.Background(),
				in: &pb.GetProductRequest{
					Id: ERROR_PRODUCT_TRIGGER,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.GetProduct(tt.args.ctx, tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("GRPCServer.GetProduct() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GRPCServer.GetProduct() = %v, want %v", got, tt.want)
			}
		})
	}
}

func mockListProductsFunc(ctx context.Context) ([]*repository.Product, error) {
	return []*repository.Product{
		{
			Id:          "1",
			Name:        "Test Product",
			Description: "Test Description",
			Price:       100,
		},
	}, nil
}

func TestGRPCServer_ListProducts(t *testing.T) {
	s := NewTestGRPCServer(t)

	s.ProductRepository.ListProductsFunc = mockListProductsFunc

	type args struct {
		ctx context.Context
		in  *pb.ListProductsRequest
	}
	tests := []struct {
		name string

		args    args
		want    *pb.ListProductsResponse
		wantErr bool
	}{
		{
			name: "List Products Success",
			args: args{
				ctx: context.Background(),
				in:  &pb.ListProductsRequest{},
			},
			want: &pb.ListProductsResponse{
				Products: []*pb.Product{
					{
						Id:          "1",
						Name:        "Test Product",
						Description: "Test Description",
						Price:       100,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.ListProducts(tt.args.ctx, tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("GRPCServer.ListProducts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GRPCServer.ListProducts() = %v, want %v", got, tt.want)
			}
		})
	}
}

func mockUpdateProductFunc(ctx context.Context, id string, update *repository.ProductUpdate) (*repository.Product, error) {
	if id == ERROR_PRODUCT_TRIGGER {
		return nil, errors.New("intentional error")
	}
	return &repository.Product{
		Id:          "1",
		Name:        "Test Product",
		Description: "Test Description",
		Price:       100,
	}, nil
}

func TestGRPCServer_UpdateProduct(t *testing.T) {
	s := NewTestGRPCServer(t)

	s.ProductRepository.UpdateProductFunc = mockUpdateProductFunc

	type args struct {
		ctx context.Context
		in  *pb.UpdateProductRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *pb.UpdateProductResponse
		wantErr bool
	}{
		{
			name: "Update Product Success",
			args: args{
				ctx: context.Background(),
				in: &pb.UpdateProductRequest{
					Id: "1",
					Update: &pb.ProductUpdate{
						Name:        "Updated Test Product",
						Description: "Updated Test Description",
						Price:       200,
					},
				},
			},
			want: &pb.UpdateProductResponse{
				Id:          "1",
				Name:        "Test Product",
				Description: "Test Description",
				Price:       100,
			},
			wantErr: false,
		},
		{
			name: "Update Product Error",
			args: args{
				ctx: context.Background(),
				in: &pb.UpdateProductRequest{
					Id: ERROR_PRODUCT_TRIGGER,
					Update: &pb.ProductUpdate{
						Name:        "Updated Test Product",
						Description: "Updated Test Description",
						Price:       200,
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := s.UpdateProduct(tt.args.ctx, tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("GRPCServer.UpdateProduct() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GRPCServer.UpdateProduct() = %v, want %v", got, tt.want)
			}
		})
	}
}

func mockDeleteProductFunc(ctx context.Context, id string) error {
	if id == ERROR_PRODUCT_TRIGGER {
		return errors.New("intentional error")
	}
	return nil
}

func TestGRPCServer_DeleteProduct(t *testing.T) {

	s := NewTestGRPCServer(t)

	s.ProductRepository.DeleteProductFunc = mockDeleteProductFunc

	type args struct {
		ctx context.Context
		in  *pb.DeleteProductRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *pb.DeleteProductResponse
		wantErr bool
	}{
		{
			name: "Delete Product Success",
			args: args{
				ctx: context.Background(),
				in: &pb.DeleteProductRequest{
					Id: "1",
				},
			},
			want:    &pb.DeleteProductResponse{},
			wantErr: false,
		},
		{
			name: "Delete Product Error",
			args: args{
				ctx: context.Background(),
				in: &pb.DeleteProductRequest{
					Id: ERROR_PRODUCT_TRIGGER,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.DeleteProduct(tt.args.ctx, tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("GRPCServer.DeleteProduct() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GRPCServer.DeleteProduct() = %v, want %v", got, tt.want)
			}
		})
	}
}
