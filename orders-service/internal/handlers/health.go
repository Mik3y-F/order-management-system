package handlers

import (
	"context"
	"log"

	pb "github.com/Mik3y-F/order-management-system/orders/api/generated"
)

func (s *GRPCServer) HealthCheck(ctx context.Context, in *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {

	log.Printf("Received: Health check request")

	return &pb.HealthCheckResponse{
		Status: "OK",
	}, nil
}
