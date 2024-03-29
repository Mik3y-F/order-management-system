package handlers

import (
	"context"
	"fmt"
	"log"

	pb "github.com/Mik3y-F/order-management-system/orders/api/generated"
	"github.com/Mik3y-F/order-management-system/orders/internal/repository"
	"github.com/Mik3y-F/order-management-system/orders/pkg"
)

func (s *GRPCServer) CreateCustomer(
	ctx context.Context, in *pb.CreateCustomerRequest) (*pb.CreateCustomerResponse, error) {

	log.Printf("Received: %v", in.GetFirstName())

	p, err := s.CustomerRepository.CreateCustomer(ctx, &repository.Customer{
		FirstName: in.GetFirstName(),
		LastName:  in.GetLastName(),
		Email:     in.GetEmail(),
		Phone:     in.GetPhone(),
	})
	if err != nil {
		return nil, Error(fmt.Errorf("failed to create customer: %w", err))
	}

	return &pb.CreateCustomerResponse{
		Id: p.Id,
	}, nil
}

func (s *GRPCServer) GetCustomer(
	ctx context.Context, in *pb.GetCustomerRequest) (*pb.GetCustomerResponse, error) {

	log.Printf("Received: %v", in.GetId())

	p, err := s.CustomerRepository.GetCustomer(ctx, in.GetId())
	if err != nil {
		return nil, fmt.Errorf("failed to get customer: %w", err)
	}

	return &pb.GetCustomerResponse{
		Id:        p.Id,
		FirstName: p.FirstName,
		LastName:  p.LastName,
		Phone:     p.Phone,
		Email:     p.Email,
	}, nil
}

func (s *GRPCServer) ListCustomers(
	ctx context.Context, in *pb.ListCustomersRequest) (*pb.ListCustomersResponse, error) {

	log.Printf("Received: %v", in)

	customers, err := s.CustomerRepository.ListCustomers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list customers: %w", err)
	}

	var responseCustomers []*pb.Customer
	for _, p := range customers {
		responseCustomers = append(responseCustomers, &pb.Customer{
			Id:        p.Id,
			FirstName: p.FirstName,
			LastName:  p.LastName,
			Phone:     p.Phone,
			Email:     p.Email,
		})
	}

	return &pb.ListCustomersResponse{
		Customers: responseCustomers,
	}, nil
}

func (s *GRPCServer) UpdateCustomer(
	ctx context.Context, in *pb.UpdateCustomerRequest) (*pb.UpdateCustomerResponse, error) {

	log.Printf("Received: %v", in)

	p, err := s.CustomerRepository.UpdateCustomer(ctx, in.GetId(), &repository.CustomerUpdate{
		FirstName: pkg.StringPtr(in.GetUpdate().GetFirstName()),
		LastName:  pkg.StringPtr(in.GetUpdate().GetLastName()),
		Phone:     pkg.StringPtr(in.GetUpdate().GetPhone()),
		Email:     pkg.StringPtr(in.GetUpdate().GetEmail()),
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

func (s *GRPCServer) DeleteCustomer(
	ctx context.Context, in *pb.DeleteCustomerRequest) (*pb.DeleteCustomerResponse, error) {

	log.Printf("Received: %v", in.GetId())

	err := s.CustomerRepository.DeleteCustomer(ctx, in.GetId())
	if err != nil {
		return nil, fmt.Errorf("failed to delete customer: %w", err)
	}

	return &pb.DeleteCustomerResponse{}, nil
}
