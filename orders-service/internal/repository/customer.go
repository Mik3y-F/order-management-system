package repository

import (
	"context"

	"github.com/Mik3y-F/order-management-system/orders/internal/service"
)

type Customer struct {
	Id        string
	FirstName string
	LastName  string
	Email     string
	Phone     string
	CreatedAt string
	UpdatedAt string
}

func (c *Customer) Validate() error {
	if c.FirstName == "" {
		return service.Errorf(service.INVALID_ERROR, "first_name is required")
	}

	if c.LastName == "" {
		return service.Errorf(service.INVALID_ERROR, "last_name is required")
	}

	if c.Email == "" {
		return service.Errorf(service.INVALID_ERROR, "email is required")
	}

	if c.Phone == "" {
		return service.Errorf(service.INVALID_ERROR, "phone is required")
	}

	return nil
}

type CustomerUpdate struct {
	FirstName *string
	LastName  *string
	Phone     *string
	Email     *string
}

type CustomerRepository interface {
	CreateCustomer(ctx context.Context, customer *Customer) (*Customer, error)
	GetCustomer(ctx context.Context, id string) (*Customer, error)
	ListCustomers(ctx context.Context) ([]*Customer, error)
	UpdateCustomer(ctx context.Context, id string, update *CustomerUpdate) (*Customer, error)
	DeleteCustomer(ctx context.Context, id string) error
}
