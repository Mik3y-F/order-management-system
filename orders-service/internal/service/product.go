package service

import "context"

type Product struct {
	Id          string
	Name        string
	Description string
	Price       uint
	CreatedAt   string
	UpdatedAt   string
}

type ProductUpdate struct {
	Name        string
	Description string
	Price       uint
}

type ProductService interface {
	CreateProduct(ctx context.Context, product *Product) (*Product, error)
	GetProduct(ctx context.Context, id string) (*Product, error)
	ListProducts(ctx context.Context) ([]*Product, error)
	UpdateProduct(ctx context.Context, id string, update *ProductUpdate) (*Product, error)
	DeleteProduct(ctx context.Context, id string) error
}

func (p *Product) Validate() error {
	if p.Name == "" {
		return Errorf(INVALID_ERROR, "name is required")
	}

	return nil
}
