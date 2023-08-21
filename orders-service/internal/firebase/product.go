package firebase

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/Mik3y-F/order-management-system/orders/internal/service"
	"google.golang.org/api/iterator"
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

	// Set CreatedAt and UpdatedAt to the current time
	currentTime := time.Now()

	product.CreatedAt = currentTime.Format(time.RFC3339)
	product.UpdatedAt = currentTime.Format(time.RFC3339)

	productModel := s.marshallProduct(product)

	docRef, _, writeErr := s.productCollection().Add(ctx, productModel)
	if writeErr != nil {
		return nil, writeErr
	}

	product.Id = docRef.ID

	return product, nil
}

func (s *ProductService) GetProduct(ctx context.Context, id string) (*service.Product, error) {
	s.CheckPreconditions()

	docRef, getErr := s.productCollection().Doc(id).Get(ctx)
	if getErr != nil {
		return nil, getErr
	}

	productModel := &ProductModel{}
	if err := docRef.DataTo(productModel); err != nil {
		return nil, err
	}

	product := s.unmarshallProduct(productModel)

	product.Id = docRef.Ref.ID

	return product, nil
}

func (s *ProductService) ListProducts(ctx context.Context) ([]*service.Product, error) {
	s.CheckPreconditions()

	iter := s.productCollection().Documents(ctx)

	var products []*service.Product

	for {
		docRef, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		productModel := &ProductModel{}
		if err := docRef.DataTo(productModel); err != nil {
			return nil, err
		}

		product := s.unmarshallProduct(productModel)
		product.Id = docRef.Ref.ID

		products = append(products, product)
	}

	return products, nil
}

func (s *ProductService) UpdateProduct(ctx context.Context, id string, update *service.ProductUpdate,
) (*service.Product, error) {
	s.CheckPreconditions()

	product, getErr := s.GetProduct(ctx, id)
	if getErr != nil {
		return nil, getErr
	}

	product.Name = update.Name
	product.Description = update.Description
	product.Price = update.Price

	timeNow := time.Now()
	product.UpdatedAt = timeNow.Format(time.RFC3339)

	productModel := s.marshallProduct(product)

	_, writeErr := s.productCollection().Doc(id).Set(ctx, productModel)

	return product, writeErr
}

func (s *ProductService) DeleteProduct(ctx context.Context, id string) error {
	s.CheckPreconditions()

	_, deleteErr := s.productCollection().Doc(id).Delete(ctx)

	return deleteErr
}

func (s *ProductService) marshallProduct(product *service.Product) *ProductModel {
	return &ProductModel{
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}
}

func (s *ProductService) unmarshallProduct(productModel *ProductModel) *service.Product {
	return &service.Product{
		Name:        productModel.Name,
		Description: productModel.Description,
		Price:       productModel.Price,
		CreatedAt:   productModel.CreatedAt,
		UpdatedAt:   productModel.UpdatedAt,
	}
}
