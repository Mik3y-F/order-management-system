package main

import (
	"context"
	"log"
	"os"

	"github.com/Mik3y-F/order-management-system/payments/internal/handlers"
)

func main() {

	ctx := context.Background()

	log.Printf("Starting server")

	bindAddress := os.Getenv("BIND_ADDRESS")
	if bindAddress == "" {
		bindAddress = "localhost"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "50051"
	}

	s := handlers.NewGRPCServer()
	if err := s.Run(ctx, bindAddress, port); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
