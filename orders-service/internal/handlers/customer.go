package handlers

import (
	"context"
	"fmt"
	"log"

	pb "github.com/Mik3y-F/order-management-system/orders/api/generated"
	"github.com/Mik3y-F/order-management-system/orders/internal/service"
)

func (s *GRPCServer) CreateCustomer(ctx context.Context, in *pb.CreateCustomerRequest) (*pb.CreateCustomerResponse, error) {

	log.Printf("Received: %v", in.GetFirstName())

	p, err := s.CustomerService.CreateCustomer(ctx, &service.Customer{
		FirstName: in.GetFirstName(),
		LastName:  in.GetLastName(),
		Email:     in.GetEmail(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create customer: %w", err)
	}

	return &pb.CreateCustomerResponse{
		Id: p.Id,
	}, nil
}

func (s *GRPCServer) GetCustomer(ctx context.Context, in *pb.GetCustomerRequest) (*pb.GetCustomerResponse, error) {

	log.Printf("Received: %v", in.GetId())

	p, err := s.CustomerService.GetCustomer(ctx, in.GetId())
	if err != nil {
		return nil, fmt.Errorf("failed to get customer: %w", err)
	}

	return &pb.GetCustomerResponse{
		Id:        p.Id,
		FirstName: p.FirstName,
		LastName:  p.LastName,
		Email:     p.Email,
	}, nil
}

func (s *GRPCServer) ListCustomers(ctx context.Context, in *pb.ListCustomersRequest) (*pb.ListCustomersResponse, error) {

	log.Printf("Received: %v", in)

	customers, err := s.CustomerService.ListCustomers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list customers: %w", err)
	}

	var responseCustomers []*pb.Customer
	for _, p := range customers {
		responseCustomers = append(responseCustomers, &pb.Customer{
			Id:        p.Id,
			FirstName: p.FirstName,
			LastName:  p.LastName,
			Email:     p.Email,
		})
	}

	return &pb.ListCustomersResponse{
		Customers: responseCustomers,
	}, nil
}

func (s *GRPCServer) UpdateCustomer(ctx context.Context, in *pb.UpdateCustomerRequest) (*pb.UpdateCustomerResponse, error) {

	log.Printf("Received: %v", in)

	p, err := s.CustomerService.UpdateCustomer(ctx, in.GetId(), &service.CustomerUpdate{
		FirstName: in.GetUpdate().GetFirstName(),
		LastName:  in.GetUpdate().GetLastName(),
		Email:     in.GetUpdate().GetEmail(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update customer: %w", err)
	}

	return &pb.UpdateCustomerResponse{
		Id:        p.Id,
		FirstName: p.FirstName,
		LastName:  p.LastName,
		Email:     p.Email,
	}, nil
}

func (s *GRPCServer) DeleteCustomer(ctx context.Context, in *pb.DeleteCustomerRequest) (*pb.DeleteCustomerResponse, error) {

	log.Printf("Received: %v", in.GetId())

	err := s.CustomerService.DeleteCustomer(ctx, in.GetId())
	if err != nil {
		return nil, fmt.Errorf("failed to delete customer: %w", err)
	}

	return &pb.DeleteCustomerResponse{}, nil
}
