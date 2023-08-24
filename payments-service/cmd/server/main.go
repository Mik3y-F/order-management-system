package main

import (
	"context"
	"log"
	"os"

	orders "github.com/Mik3y-F/order-management-system/orders/pkg/client"
	db "github.com/Mik3y-F/order-management-system/payments/internal/firebase"
	ecom_grpc "github.com/Mik3y-F/order-management-system/payments/internal/handlers/grpc"
	"github.com/Mik3y-F/order-management-system/payments/internal/mpesa"
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

	s := ecom_grpc.NewGRPCServer()

	mpesaService := mpesa.NewMpesaService()

	// Setup order service client
	conn, err := orders.ConnectToOrderService("localhost:50051")
	if err != nil {
		log.Fatalf("Failed to connect to order service: %v", err)
	}
	orderClient := orders.NewGrpcOrderClient(conn)

	// Setup firebase client and firestore service
	firebase := db.NewFirebaseService()
	firestoreClient, err := firebase.GetApp().Firestore(ctx)
	if err != nil {
		log.Fatalf("failed to create firestore client: %v", err)
	}
	defer firestoreClient.Close()

	firestoreService := db.NewFirestoreService(firestoreClient)
	paymentRepository := db.NewPaymentsRepository(firestoreService)

	paymentService := mpesa.NewPaymentsService(mpesaService, orderClient, paymentRepository)

	// Register internal services
	s.PaymentsService = paymentService

	if err := s.Run(ctx, bindAddress, port); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
