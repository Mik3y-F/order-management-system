package main

import (
	"context"
	"log"
	"os"

	"github.com/Mik3y-F/order-management-system/payments/internal/handlers"
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

	s := handlers.NewGRPCServer()

	mpesaService := mpesa.NewMpesaService()
	paymentService := mpesa.NewPaymentsService(mpesaService)

	// Register internal services
	s.PaymentsService = paymentService

	if err := s.Run(ctx, bindAddress, port); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
