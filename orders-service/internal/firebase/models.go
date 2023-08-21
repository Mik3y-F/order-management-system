package firebase

type ProductModel struct {
	Name        string  `firestore:"name"`
	Description string  `firestore:"description"`
	Price       float32 `firestore:"price"`
	CreatedAt   string  `firestore:"created_at"`
	UpdatedAt   string  `firestore:"updated_at"`
}

type CustomerModel struct {
	FirstName string `firestore:"first_name"`
	LastName  string `firestore:"last_name"`
	Email     string `firestore:"email"`
	CreatedAt string `firestore:"created_at"`
	UpdatedAt string `firestore:"updated_at"`
}
