package mock

import (
	"context"

	"github.com/Mik3y-F/order-management-system/orders/internal/service"
)

var _ service.ProductService = (*ProductService)(nil)

type ProductService struct {
	CreateProductFunc func(ctx context.Context, p *service.Product) (*service.Product, error)
}

func (m *ProductService) CreateProduct(ctx context.Context, p *service.Product) (*service.Product, error) {
	return m.CreateProductFunc(ctx, p)
}
