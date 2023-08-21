package mock

import (
	"context"

	"github.com/Mik3y-F/order-management-system/orders/internal/service"
)

var _ service.CustomerService = (*CustomerService)(nil)

type CustomerService struct {
	CreateCustomerFunc func(ctx context.Context, p *service.Customer) (*service.Customer, error)
	GetCustomerFunc    func(ctx context.Context, id string) (*service.Customer, error)
	ListCustomersFunc  func(ctx context.Context) ([]*service.Customer, error)
	UpdateCustomerFunc func(ctx context.Context, id string, update *service.CustomerUpdate) (*service.Customer, error)
	DeleteCustomerFunc func(ctx context.Context, id string) error
}

func (m *CustomerService) CreateCustomer(ctx context.Context, p *service.Customer) (*service.Customer, error) {
	return m.CreateCustomerFunc(ctx, p)
}

func (m *CustomerService) GetCustomer(ctx context.Context, id string) (*service.Customer, error) {
	return m.GetCustomerFunc(ctx, id)
}

func (m *CustomerService) ListCustomers(ctx context.Context) ([]*service.Customer, error) {
	return m.ListCustomersFunc(ctx)
}

func (m *CustomerService) UpdateCustomer(ctx context.Context, id string, update *service.CustomerUpdate,
) (*service.Customer, error) {
	return m.UpdateCustomerFunc(ctx, id, update)
}

func (m *CustomerService) DeleteCustomer(ctx context.Context, id string) error {
	return m.DeleteCustomerFunc(ctx, id)
}
