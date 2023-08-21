package firebase

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/Mik3y-F/order-management-system/orders/internal/service"
	"google.golang.org/api/iterator"
)

var _ service.CustomerService = (*CustomerService)(nil)

type CustomerService struct {
	db *FirestoreService
}

func NewCustomerService(db *FirestoreService) *CustomerService {
	return &CustomerService{
		db: db,
	}
}

func (s *CustomerService) CheckPreconditions() {
	if s.db == nil {
		panic("no DB service provided")
	}
}

func (s *CustomerService) customerCollection() *firestore.CollectionRef {
	s.CheckPreconditions()

	return s.db.client.Collection("customers")
}

func (s *CustomerService) CreateCustomer(ctx context.Context, customer *service.Customer) (*service.Customer, error) {
	s.CheckPreconditions()

	currentTime := time.Now()
	customer.CreatedAt = currentTime.Format(time.RFC3339)
	customer.UpdatedAt = currentTime.Format(time.RFC3339)

	customerModel := s.marshallCustomer(customer)

	docRef, _, writeErr := s.customerCollection().Add(ctx, customerModel)
	if writeErr != nil {
		return nil, writeErr
	}

	customer.Id = docRef.ID

	return customer, nil
}

func (s *CustomerService) GetCustomer(ctx context.Context, id string) (*service.Customer, error) {

	s.CheckPreconditions()

	docRef, getErr := s.customerCollection().Doc(id).Get(ctx)
	if getErr != nil {
		return nil, getErr
	}

	customerModel := &CustomerModel{}
	if err := docRef.DataTo(customerModel); err != nil {
		return nil, err
	}

	customer := s.unmarshallCustomer(customerModel)

	customer.Id = docRef.Ref.ID

	return customer, nil
}

func (s *CustomerService) ListCustomers(ctx context.Context) ([]*service.Customer, error) {
	s.CheckPreconditions()

	iter := s.customerCollection().Documents(ctx)

	var customers []*service.Customer

	for {
		docRef, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		customerModel := &CustomerModel{}
		if err := docRef.DataTo(customerModel); err != nil {
			return nil, err
		}

		customer := s.unmarshallCustomer(customerModel)
		customer.Id = docRef.Ref.ID

		customers = append(customers, customer)
	}

	return customers, nil
}

func (s *CustomerService) UpdateCustomer(
	ctx context.Context, id string, update *service.CustomerUpdate) (*service.Customer, error) {

	s.CheckPreconditions()

	customer, getErr := s.GetCustomer(ctx, id)
	if getErr != nil {
		return nil, getErr
	}

	// todo: update only the fields that are not empty
	customer.FirstName = update.FirstName
	customer.LastName = update.LastName
	customer.Email = update.Email

	timeNow := time.Now()
	customer.UpdatedAt = timeNow.Format(time.RFC3339)

	customerModel := s.marshallCustomer(customer)

	_, writeErr := s.customerCollection().Doc(id).Set(ctx, customerModel)

	return customer, writeErr

}

func (s *CustomerService) DeleteCustomer(ctx context.Context, id string) error {

	_, err := s.customerCollection().Doc(id).Delete(ctx)

	return err
}

func (s *CustomerService) marshallCustomer(customer *service.Customer) *CustomerModel {

	return &CustomerModel{
		FirstName: customer.FirstName,
		LastName:  customer.LastName,
		Email:     customer.Email,
		CreatedAt: customer.CreatedAt,
		UpdatedAt: customer.UpdatedAt,
	}
}

func (s *CustomerService) unmarshallCustomer(customerModel *CustomerModel) *service.Customer {

	return &service.Customer{
		FirstName: customerModel.FirstName,
		LastName:  customerModel.LastName,
		Email:     customerModel.Email,
		CreatedAt: customerModel.CreatedAt,
		UpdatedAt: customerModel.UpdatedAt,
	}
}
