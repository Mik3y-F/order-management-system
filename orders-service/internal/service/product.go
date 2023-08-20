package service

import "context"

type Product struct {
	Id          string
	Name        string
	Description string
	Price       float32
	Stock       int32
	CreatedAt   string
	UpdatedAt   string
}

type ProductService interface {
	CreateProduct(ctx context.Context, product *Product) (*Product, error)
}
