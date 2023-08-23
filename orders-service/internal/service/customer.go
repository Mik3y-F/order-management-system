package service

import "context"

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
		return Errorf(INVALID_ERROR, "first_name is required")
	}

	if c.LastName == "" {
		return Errorf(INVALID_ERROR, "last_name is required")
	}

	if c.Email == "" {
		return Errorf(INVALID_ERROR, "email is required")
	}

	if c.Phone == "" {
		return Errorf(INVALID_ERROR, "phone is required")
	}

	return nil
}

type CustomerUpdate struct {
	FirstName *string
	LastName  *string
	Phone     *string
	Email     *string
}

type CustomerService interface {
	CreateCustomer(ctx context.Context, customer *Customer) (*Customer, error)
	GetCustomer(ctx context.Context, id string) (*Customer, error)
	ListCustomers(ctx context.Context) ([]*Customer, error)
	UpdateCustomer(ctx context.Context, id string, update *CustomerUpdate) (*Customer, error)
	DeleteCustomer(ctx context.Context, id string) error
}
