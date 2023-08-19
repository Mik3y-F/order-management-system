package handlers_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	pb "github.com/Mik3y-F/order-management-system/orders/api/generated"
	"github.com/Mik3y-F/order-management-system/orders/internal/service"
)

const (
	ERROR_PRODUCT_TRIGGER = "Error Product"
)

func mockCreateProductFunc(ctx context.Context, p *service.Product) (*service.Product, error) {
	if p.Name == ERROR_PRODUCT_TRIGGER {
		return nil, errors.New("intentional error")
	}
	return &service.Product{
		Id:          "1",
		Name:        "Test Product",
		Description: "Test Description",
		Price:       100,
	}, nil
}

func TestGRPCServer_CreateProduct(t *testing.T) {

	s := NewTestGRPCServer(t)

	s.ProductService.CreateProductFunc = mockCreateProductFunc

	type args struct {
		ctx     context.Context
		request *pb.CreateProductRequest
	}
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
