package main

import (
	"log"

	"github.com/Mik3y-F/order-management-system/orders/internal/handlers"
)

func main() {

	log.Printf("Starting server")

	s := handlers.NewGRPCServer()

	if err := s.Run(); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
