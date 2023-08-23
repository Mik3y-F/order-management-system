package main

import (
	"context"
	"log"
	"os"

	db "github.com/Mik3y-F/order-management-system/orders/internal/firebase"
	"github.com/Mik3y-F/order-management-system/orders/internal/handlers"
)

const (
	BIND_ADDRESS = "BIND_ADDRESS"
	PORT         = "PORT"

	DEFAULT_BIND_ADDRESS = "localhost"
	DEFAULT_PORT         = "50051"
)

func main() {

	ctx := context.Background()

	log.Printf("Starting server")

	bindAddress := os.Getenv(BIND_ADDRESS)
	if bindAddress == "" {
		bindAddress = DEFAULT_BIND_ADDRESS
	}

	port := os.Getenv(PORT)
	if port == "" {
		port = DEFAULT_PORT
	}

	s := handlers.NewGRPCServer()

	firebase := db.NewFirebaseService()
	firestoreClient, err := firebase.GetApp().Firestore(ctx)
	if err != nil {
		log.Fatalf("failed to create firestore client: %v", err)
	}
	defer firestoreClient.Close()

	firestoreService := db.NewFirestoreService(firestoreClient)

	productService := db.NewProductService(firestoreService)
	customerService := db.NewCustomerService(firestoreService)
	orderService := db.NewOrderService(firestoreService)

	s.ProductService = productService
	s.CustomerService = customerService
	s.OrderService = orderService

	if err := s.Run(ctx, bindAddress, port); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
