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

var _ repository.CustomerRepository = (*CustomerRepository)(nil)

type CustomerRepository struct {
	db *FirestoreService
}

func NewCustomerService(db *FirestoreService) *CustomerRepository {
	return &CustomerRepository{
		db: db,
	}
}

func (s *CustomerRepository) CheckPreconditions() {
	if s.db == nil {
		panic("no DB service provided")
	}
}

func (s *CustomerRepository) customerCollection() *firestore.CollectionRef {
	s.CheckPreconditions()

	return s.db.client.Collection("customers")
}

func (s *CustomerRepository) CreateCustomer(ctx context.Context, customer *repository.Customer) (*repository.Customer, error) {
	s.CheckPreconditions()

	currentTime := time.Now()
	customer.CreatedAt = currentTime.Format(time.RFC3339)
	customer.UpdatedAt = currentTime.Format(time.RFC3339)

	err := customer.Validate()
	if err != nil {
		return nil, service.Errorf(service.INVALID_ERROR, "invalid customer provided: %v", err)
	}
	customerModel := s.marshallCustomer(customer)

	docRef, _, err := s.customerCollection().Add(ctx, customerModel)
	if err != nil {
		return nil, service.Errorf(service.INTERNAL_ERROR, "failed to create customer: %v", err)
	}

	customer.Id = docRef.ID

	return customer, nil
}

func (s *CustomerRepository) GetCustomer(ctx context.Context, id string) (*repository.Customer, error) {

	s.CheckPreconditions()

	if id == "" {
		return nil, service.Errorf(service.INVALID_ERROR, "id is required")
	}

	docRef, err := s.customerCollection().Doc(id).Get(ctx)
	if status.Code(err) == codes.NotFound {
		return nil, service.Errorf(service.NOT_FOUND_ERROR, "customer not found")
	} else if err != nil {
		return nil, service.Errorf(service.INTERNAL_ERROR, "failed to get customer: %v", err)
	}

	customerModel := &CustomerModel{}
	if err := docRef.DataTo(customerModel); err != nil {
		return nil, service.Errorf(service.INTERNAL_ERROR, "failed to unmarshall customer: %v", err)
	}

	customer := s.unmarshallCustomer(customerModel)

	customer.Id = docRef.Ref.ID

	return customer, nil
}

func (s *CustomerRepository) ListCustomers(ctx context.Context) ([]*repository.Customer, error) {
	s.CheckPreconditions()

	iter := s.customerCollection().Documents(ctx)

	var customers []*repository.Customer

	for {
		docRef, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, service.Errorf(service.INTERNAL_ERROR, "failed to iterate customers: %v", err)
		}

		customerModel := &CustomerModel{}
		if err := docRef.DataTo(customerModel); err != nil {
			return nil, service.Errorf(service.INTERNAL_ERROR, "failed to unmarshall customer: %v", err)
		}

		customer := s.unmarshallCustomer(customerModel)
		customer.Id = docRef.Ref.ID

		customers = append(customers, customer)
	}

	return customers, nil
}

func (s *CustomerRepository) UpdateCustomer(
	ctx context.Context, id string, update *repository.CustomerUpdate) (*repository.Customer, error) {

	s.CheckPreconditions()

	customer, err := s.GetCustomer(ctx, id)
	if err != nil {
		return nil, err
	}

	if c := update.FirstName; c != nil {
		customer.FirstName = *c
	}

	if c := update.LastName; c != nil {
		customer.LastName = *c
	}

	if c := update.Email; c != nil {
		customer.Email = *c
	}

	if c := update.Phone; c != nil {
		customer.Phone = *c
	}

	timeNow := time.Now()
	customer.UpdatedAt = timeNow.Format(time.RFC3339)

	err = customer.Validate()
	if err != nil {
		return nil, service.Errorf(service.INVALID_ERROR, "invalid customer details provided: %v", err)
	}

	customerModel := s.marshallCustomer(customer)

	_, err = s.customerCollection().Doc(id).Set(ctx, customerModel)
	if err != nil {
		return nil, service.Errorf(service.INTERNAL_ERROR, "failed to update customer: %v", err)
	}

	return customer, nil

}

func (s *CustomerRepository) DeleteCustomer(ctx context.Context, id string) error {
	s.CheckPreconditions()

	if id == "" {
		return service.Errorf(service.INVALID_ERROR, "id is required")
	}

	_, err := s.customerCollection().Doc(id).Delete(ctx)

	return err
}

func (s *CustomerRepository) marshallCustomer(customer *repository.Customer) *CustomerModel {

	return &CustomerModel{
		FirstName: customer.FirstName,
		LastName:  customer.LastName,
		Email:     customer.Email,
		Phone:     customer.Phone,
		CreatedAt: customer.CreatedAt,
		UpdatedAt: customer.UpdatedAt,
	}
}

func (s *CustomerRepository) unmarshallCustomer(customerModel *CustomerModel) *repository.Customer {

	return &repository.Customer{
		FirstName: customerModel.FirstName,
		LastName:  customerModel.LastName,
		Email:     customerModel.Email,
		Phone:     customerModel.Phone,
		CreatedAt: customerModel.CreatedAt,
		UpdatedAt: customerModel.UpdatedAt,
	}
}
