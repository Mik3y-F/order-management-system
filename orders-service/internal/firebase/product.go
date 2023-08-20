package firebase

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/Mik3y-F/order-management-system/orders/internal/service"
)

var _ service.ProductService = (*ProductService)(nil)

type ProductService struct {
	db *FirestoreService
}

func NewProductService(db *FirestoreService) *ProductService {
	return &ProductService{
		db: db,
	}
}

func (s *ProductService) CheckPreconditions() {
	if s.db == nil {
		panic("no DB service provided")
	}
}

func (s *ProductService) productCollection() *firestore.CollectionRef {
	s.CheckPreconditions()

	return s.db.client.Collection("products")
}

func (s *ProductService) CreateProduct(ctx context.Context, product *service.Product) (*service.Product, error) {
	s.CheckPreconditions()

	productModel := s.marshallProduct(product)

	docRef, _, writeErr := s.productCollection().Add(ctx, productModel)
	if writeErr != nil {
		return nil, writeErr
	}

	product.Id = docRef.ID

	return product, nil
}

func (s *ProductService) marshallProduct(product *service.Product) *ProductModel {
	return &ProductModel{
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
	}
}
