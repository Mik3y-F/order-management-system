package firebase

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/Mik3y-F/order-management-system/orders/internal/repository"
	"github.com/Mik3y-F/order-management-system/orders/internal/service"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ repository.ProductRepository = (*ProductRepository)(nil)

type ProductRepository struct {
	db *FirestoreService
}

func NewProductService(db *FirestoreService) *ProductRepository {
	return &ProductRepository{
		db: db,
	}
}

func (r *ProductRepository) CheckPreconditions() {
	if r.db == nil {
		panic("no DB service provided")
	}
}

func (r *ProductRepository) productCollection() *firestore.CollectionRef {
	r.CheckPreconditions()

	return r.db.client.Collection("products")
}

func (r *ProductRepository) CreateProduct(ctx context.Context, product *repository.Product) (*repository.Product, error) {
	r.CheckPreconditions()

	// Set CreatedAt and UpdatedAt to the current time
	currentTime := time.Now()

	product.CreatedAt = currentTime.Format(time.RFC3339)
	product.UpdatedAt = currentTime.Format(time.RFC3339)

	err := product.Validate()
	if err != nil {
		return nil, service.Errorf(service.INVALID_ERROR, "invalid product provided: %v", err)
	}

	productModel := r.marshallProduct(product)

	docRef, _, writeErr := r.productCollection().Add(ctx, productModel)
	if writeErr != nil {
		return nil, writeErr
	}

	product.Id = docRef.ID

	return product, nil
}

func (r *ProductRepository) GetProduct(ctx context.Context, id string) (*repository.Product, error) {
	r.CheckPreconditions()

	if id == "" {
		return nil, service.Errorf(service.INVALID_ERROR, "id is required")
	}

	docRef, getErr := r.productCollection().Doc(id).Get(ctx)
	if status.Code(getErr) == codes.NotFound {
		return nil, service.Errorf(service.NOT_FOUND_ERROR, "product not found")
	} else if getErr != nil {
		return nil, service.Errorf(service.INTERNAL_ERROR, "failed to get product: %v", getErr)
	}

	productModel := &ProductModel{}
	if err := docRef.DataTo(productModel); err != nil {
		return nil, service.Errorf(service.INTERNAL_ERROR, "failed to unmarshall product: %v", err)
	}

	product := r.unmarshallProduct(productModel)

	product.Id = docRef.Ref.ID

	return product, nil
}

func (r *ProductRepository) ListProducts(ctx context.Context) ([]*repository.Product, error) {
	r.CheckPreconditions()

	iter := r.productCollection().Documents(ctx)

	var products []*repository.Product

	for {
		docRef, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, service.Errorf(service.INTERNAL_ERROR, "failed to iterate products: %v", err)
		}

		productModel := &ProductModel{}
		if err := docRef.DataTo(productModel); err != nil {
			return nil, service.Errorf(service.INTERNAL_ERROR, "failed to unmarshall product: %v", err)
		}

		product := r.unmarshallProduct(productModel)
		product.Id = docRef.Ref.ID

		products = append(products, product)
	}

	return products, nil
}

func (r *ProductRepository) UpdateProduct(ctx context.Context, id string, update *repository.ProductUpdate,
) (*repository.Product, error) {
	r.CheckPreconditions()

	product, getErr := r.GetProduct(ctx, id)
	if getErr != nil {
		return nil, getErr
	}

	if p := update.Name; p != nil {
		product.Name = *p
	}

	if p := update.Description; p != nil {
		product.Description = *p
	}

	if p := update.Price; p != nil {
		product.Price = *p
	}

	err := product.Validate()
	if err != nil {
		return nil, service.Errorf(service.INVALID_ERROR, "invalid product details provided: %v", err)
	}

	timeNow := time.Now()
	product.UpdatedAt = timeNow.Format(time.RFC3339)

	productModel := r.marshallProduct(product)

	_, err = r.productCollection().Doc(id).Set(ctx, productModel)
	if err != nil {
		return nil, service.Errorf(service.INTERNAL_ERROR, "failed to update product: %v", err)
	}

	return product, nil
}

func (r *ProductRepository) DeleteProduct(ctx context.Context, id string) error {
	r.CheckPreconditions()

	_, err := r.productCollection().Doc(id).Delete(ctx)
	if err != nil {
		return service.Errorf(service.INTERNAL_ERROR, "failed to delete product: %v", err)
	}

	return err
}

func (r *ProductRepository) marshallProduct(product *repository.Product) *ProductModel {
	return &ProductModel{
		Name:        product.Name,
		Description: product.Description,
		Price:       int(product.Price),
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}
}

func (r *ProductRepository) unmarshallProduct(productModel *ProductModel) *repository.Product {
	return &repository.Product{
		Name:        productModel.Name,
		Description: productModel.Description,
		Price:       uint(productModel.Price),
		CreatedAt:   productModel.CreatedAt,
		UpdatedAt:   productModel.UpdatedAt,
	}
}
