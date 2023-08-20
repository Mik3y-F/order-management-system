package handlers

import (
	"context"
	"fmt"
	"log"

	pb "github.com/Mik3y-F/order-management-system/orders/api/generated"
	"github.com/Mik3y-F/order-management-system/orders/internal/service"
)

func (s *GRPCServer) CreateProduct(ctx context.Context, in *pb.CreateProductRequest) (*pb.CreateProductResponse, error) {

	log.Printf("Received: %v", in.GetName())

	p, err := s.ProductService.CreateProduct(ctx, &service.Product{
		Name:        in.GetName(),
		Description: in.GetDescription(),
		Price:       in.GetPrice(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	return &pb.CreateProductResponse{
		Id: p.Id,
	}, nil
}
