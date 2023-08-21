package service

import "context"

type Customer struct {
	Id        string
	FirstName string
	LastName  string
	Email     string
	CreatedAt string
	UpdatedAt string
}

type CustomerUpdate struct {
	FirstName string
	LastName  string
	Email     string
}

type CustomerService interface {
	CreateCustomer(ctx context.Context, customer *Customer) (*Customer, error)
	GetCustomer(ctx context.Context, id string) (*Customer, error)
	ListCustomers(ctx context.Context) ([]*Customer, error)
	UpdateCustomer(ctx context.Context, id string, update *CustomerUpdate) (*Customer, error)
	DeleteCustomer(ctx context.Context, id string) error
}
