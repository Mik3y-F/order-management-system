package firebase

type ProductModel struct {
	Name        string `firestore:"name"`
	Description string `firestore:"description"`
	Price       int    `firestore:"price"`
	CreatedAt   string `firestore:"created_at"`
	UpdatedAt   string `firestore:"updated_at"`
}

type CustomerModel struct {
	FirstName string `firestore:"first_name"`
	LastName  string `firestore:"last_name"`
	Email     string `firestore:"email"`
	CreatedAt string `firestore:"created_at"`
	UpdatedAt string `firestore:"updated_at"`
}

type OrderModel struct {
	CustomerId string            `firestore:"customer_id"`
	Items      []*OrderItemModel `firestore:"items"`
	CreatedAt  string            `firestore:"created_at"`
	UpdatedAt  string            `firestore:"updated_at"`
}

type OrderItemModel struct {
	ProductId string `firestore:"product_id"`
	Quantity  int    `firestore:"quantity"`
	CreatedAt string `firestore:"created_at"`
	UpdatedAt string `firestore:"updated_at"`
}
