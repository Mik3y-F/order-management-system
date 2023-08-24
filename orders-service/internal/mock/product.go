package mock

import (
	"context"

	"github.com/Mik3y-F/order-management-system/orders/internal/repository"
)

var _ repository.ProductRepository = (*ProductRepository)(nil)

type ProductRepository struct {
	CreateProductFunc func(ctx context.Context, p *repository.Product) (*repository.Product, error)
	GetProductFunc    func(ctx context.Context, id string) (*repository.Product, error)
	ListProductsFunc  func(ctx context.Context) ([]*repository.Product, error)
	UpdateProductFunc func(ctx context.Context, id string, update *repository.ProductUpdate) (*repository.Product, error)
	DeleteProductFunc func(ctx context.Context, id string) error
}

func (m *ProductRepository) CreateProduct(ctx context.Context, p *repository.Product) (*repository.Product, error) {
	return m.CreateProductFunc(ctx, p)
}

func (m *ProductRepository) GetProduct(ctx context.Context, id string) (*repository.Product, error) {
	return m.GetProductFunc(ctx, id)
}

func (m *ProductRepository) ListProducts(ctx context.Context) ([]*repository.Product, error) {
	return m.ListProductsFunc(ctx)
}

func (m *ProductRepository) UpdateProduct(ctx context.Context, id string, update *repository.ProductUpdate,
) (*repository.Product, error) {
	return m.UpdateProductFunc(ctx, id, update)
}

func (m *ProductRepository) DeleteProduct(ctx context.Context, id string) error {
	return m.DeleteProductFunc(ctx, id)
}
