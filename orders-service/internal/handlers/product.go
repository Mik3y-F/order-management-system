package handlers

import (
	"context"
	"fmt"
	"log"

	pb "github.com/Mik3y-F/order-management-system/orders/api/generated"
	"github.com/Mik3y-F/order-management-system/orders/internal/repository"
	"github.com/Mik3y-F/order-management-system/orders/pkg"
)

func (s *GRPCServer) CreateProduct(ctx context.Context, in *pb.CreateProductRequest) (*pb.CreateProductResponse, error) {

	log.Printf("Received: %v", in.GetName())

	p, err := s.ProductRepository.CreateProduct(ctx, &repository.Product{
		Name:        in.GetName(),
		Description: in.GetDescription(),
		Price:       uint(in.GetPrice()),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	return &pb.CreateProductResponse{
		Id: p.Id,
	}, nil
}

func (s *GRPCServer) GetProduct(ctx context.Context, in *pb.GetProductRequest) (*pb.GetProductResponse, error) {

	log.Printf("Received: %v", in.GetId())

	p, err := s.ProductRepository.GetProduct(ctx, in.GetId())
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	return &pb.GetProductResponse{
		Id:          p.Id,
		Name:        p.Name,
		Description: p.Description,
		Price:       uint32(p.Price),
	}, nil
}

func (s *GRPCServer) ListProducts(ctx context.Context, in *pb.ListProductsRequest) (*pb.ListProductsResponse, error) {

	log.Printf("Received: %v", in)

	products, err := s.ProductRepository.ListProducts(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list products: %w", err)
	}

	var responseProducts []*pb.Product
	for _, p := range products {
		responseProducts = append(responseProducts, &pb.Product{
			Id:          p.Id,
			Name:        p.Name,
			Description: p.Description,
			Price:       uint32(p.Price),
		})
	}

	return &pb.ListProductsResponse{
		Products: responseProducts,
	}, nil
}

func (s *GRPCServer) UpdateProduct(ctx context.Context, in *pb.UpdateProductRequest) (*pb.UpdateProductResponse, error) {

	log.Printf("Received: %v", in)

	p, err := s.ProductRepository.UpdateProduct(ctx, in.GetId(), &repository.ProductUpdate{
		Name:        pkg.StringPtr(in.GetUpdate().GetName()),
		Description: pkg.StringPtr(in.GetUpdate().GetDescription()),
		Price:       pkg.UintPtr(uint(in.GetUpdate().GetPrice())),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	return &pb.UpdateProductResponse{
		Id:          p.Id,
		Name:        p.Name,
		Description: p.Description,
		Price:       uint32(p.Price),
	}, nil
}

func (s *GRPCServer) DeleteProduct(ctx context.Context, in *pb.DeleteProductRequest) (*pb.DeleteProductResponse, error) {

	log.Printf("Received: %v", in)

	err := s.ProductRepository.DeleteProduct(ctx, in.GetId())
	if err != nil {
		return nil, fmt.Errorf("failed to delete product: %w", err)
	}

	return &pb.DeleteProductResponse{}, nil
}
