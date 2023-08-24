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
	ERROR_CUSTOMER_TRIGGER = "Error Customer"
)

func mockCreateCustomerFunc(ctx context.Context, c *repository.Customer) (*repository.Customer, error) {
	if c.FirstName == ERROR_CUSTOMER_TRIGGER {
		return nil, errors.New("intentional error")
	}
	return &repository.Customer{
		Id:        "1",
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
	}, nil
}

func TestGRPCServer_CreateCustomer(t *testing.T) {

	s := NewTestGRPCServer(t)

	s.CustomerRepository.CreateCustomerFunc = mockCreateCustomerFunc

	tests := []struct {
		name    string
		request *pb.CreateCustomerRequest
		want    *pb.CreateCustomerResponse
		wantErr bool
	}{
		{
			name: "Create Customer Success",
			request: &pb.CreateCustomerRequest{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john.doe@example.com",
			},
			want: &pb.CreateCustomerResponse{
				Id: "1",
			},
			wantErr: false,
		},
		{
			name: "Create Customer Error",
			request: &pb.CreateCustomerRequest{
				FirstName: ERROR_CUSTOMER_TRIGGER,
				LastName:  "Doe",
				Email:     "error.customer@example.com",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := s.GRPCServer.CreateCustomer(context.Background(), tt.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("GRPCServer.CreateCustomer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GRPCServer.CreateCustomer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func mockGetCustomerFunc(ctx context.Context, id string) (*repository.Customer, error) {
	if id == ERROR_CUSTOMER_TRIGGER {
		return nil, errors.New("intentional error")
	}
	return &repository.Customer{
		Id:        "1",
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
	}, nil
}

func TestGRPCServer_GetCustomer(t *testing.T) {

	s := NewTestGRPCServer(t)

	s.CustomerRepository.GetCustomerFunc = mockGetCustomerFunc

	type args struct {
		ctx context.Context
		in  *pb.GetCustomerRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *pb.GetCustomerResponse
		wantErr bool
	}{
		{
			name: "Get Customer Success",
			args: args{
				ctx: context.Background(),
				in: &pb.GetCustomerRequest{
					Id: "1",
				},
			},
			want: &pb.GetCustomerResponse{
				Id:        "1",
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john.doe@example.com",
			},
			wantErr: false,
		},
		{
			name: "Get Customer Error",
			args: args{
				ctx: context.Background(),
				in: &pb.GetCustomerRequest{
					Id: ERROR_CUSTOMER_TRIGGER,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.GetCustomer(tt.args.ctx, tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("GRPCServer.GetCustomer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GRPCServer.GetCustomer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func mockListCustomersFunc(ctx context.Context) ([]*repository.Customer, error) {
	return []*repository.Customer{
		{
			Id:        "1",
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@example.com",
		},
	}, nil
}

func TestGRPCServer_ListCustomers(t *testing.T) {
	s := NewTestGRPCServer(t)

	s.CustomerRepository.ListCustomersFunc = mockListCustomersFunc

	type args struct {
		ctx context.Context
		in  *pb.ListCustomersRequest
	}
	tests := []struct {
		name string

		args    args
		want    *pb.ListCustomersResponse
		wantErr bool
	}{
		{
			name: "List Customers Success",
			args: args{
				ctx: context.Background(),
				in:  &pb.ListCustomersRequest{},
			},
			want: &pb.ListCustomersResponse{
				Customers: []*pb.Customer{
					{
						Id:        "1",
						FirstName: "John",
						LastName:  "Doe",
						Email:     "john.doe@example.com",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.ListCustomers(tt.args.ctx, tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("GRPCServer.ListCustomers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GRPCServer.ListCustomers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func mockUpdateCustomerFunc(ctx context.Context, id string, update *repository.CustomerUpdate) (*repository.Customer, error) {
	if id == ERROR_CUSTOMER_TRIGGER {
		return nil, errors.New("intentional error")
	}
	return &repository.Customer{
		Id:        "1",
		FirstName: *update.FirstName,
		LastName:  *update.LastName,
		Email:     *update.Email,
	}, nil
}

func TestGRPCServer_UpdateCustomer(t *testing.T) {
	s := NewTestGRPCServer(t)

	s.CustomerRepository.UpdateCustomerFunc = mockUpdateCustomerFunc

	type args struct {
		ctx context.Context
		in  *pb.UpdateCustomerRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *pb.UpdateCustomerResponse
		wantErr bool
	}{
		{
			name: "Update Customer Success",
			args: args{
				ctx: context.Background(),
				in: &pb.UpdateCustomerRequest{
					Id: "1",
					Update: &pb.CustomerUpdate{
						FirstName: "Jane",
						LastName:  "Doe",
						Email:     "jane.doe@example.com",
					},
				},
			},
			want: &pb.UpdateCustomerResponse{
				Id:        "1",
				FirstName: "Jane",
				LastName:  "Doe",
				Email:     "jane.doe@example.com",
			},
			wantErr: false,
		},
		{
			name: "Update Customer Error",
			args: args{
				ctx: context.Background(),
				in: &pb.UpdateCustomerRequest{
					Id: ERROR_CUSTOMER_TRIGGER,
					Update: &pb.CustomerUpdate{
						FirstName: "Jane",
						LastName:  "Doe",
						Email:     "jane.doe@example.com",
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.UpdateCustomer(tt.args.ctx, tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("GRPCServer.UpdateCustomer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GRPCServer.UpdateCustomer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func mockDeleteCustomerFunc(ctx context.Context, id string) error {
	if id == ERROR_CUSTOMER_TRIGGER {
		return errors.New("intentional error")
	}
	return nil
}

func TestGRPCServer_DeleteCustomer(t *testing.T) {
	s := NewTestGRPCServer(t)

	s.CustomerRepository.DeleteCustomerFunc = mockDeleteCustomerFunc

	type args struct {
		ctx context.Context
		in  *pb.DeleteCustomerRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *pb.DeleteCustomerResponse
		wantErr bool
	}{
		{
			name: "Delete Customer Success",
			args: args{
				ctx: context.Background(),
				in: &pb.DeleteCustomerRequest{
					Id: "1",
				},
			},
			want:    &pb.DeleteCustomerResponse{},
			wantErr: false,
		},
		{
			name: "Delete Customer Error",
			args: args{
				ctx: context.Background(),
				in: &pb.DeleteCustomerRequest{
					Id: ERROR_CUSTOMER_TRIGGER,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.DeleteCustomer(tt.args.ctx, tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("GRPCServer.DeleteCustomer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GRPCServer.DeleteCustomer() = %v, want %v", got, tt.want)
			}
		})
	}
}
