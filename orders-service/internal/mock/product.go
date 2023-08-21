package mock

import (
	"context"

	"github.com/Mik3y-F/order-management-system/orders/internal/service"
)

var _ service.ProductService = (*ProductService)(nil)

type ProductService struct {
	CreateProductFunc func(ctx context.Context, p *service.Product) (*service.Product, error)
	GetProductFunc    func(ctx context.Context, id string) (*service.Product, error)
	ListProductsFunc  func(ctx context.Context) ([]*service.Product, error)
	UpdateProductFunc func(ctx context.Context, id string, update *service.ProductUpdate) (*service.Product, error)
	DeleteProductFunc func(ctx context.Context, id string) error
}

func (m *ProductService) CreateProduct(ctx context.Context, p *service.Product) (*service.Product, error) {
	return m.CreateProductFunc(ctx, p)
}

func (m *ProductService) GetProduct(ctx context.Context, id string) (*service.Product, error) {
	return m.GetProductFunc(ctx, id)
}

func (m *ProductService) ListProducts(ctx context.Context) ([]*service.Product, error) {
	return m.ListProductsFunc(ctx)
}

func (m *ProductService) UpdateProduct(ctx context.Context, id string, update *service.ProductUpdate,
) (*service.Product, error) {
	return m.UpdateProductFunc(ctx, id, update)
}

func (m *ProductService) DeleteProduct(ctx context.Context, id string) error {
	return m.DeleteProductFunc(ctx, id)
}
