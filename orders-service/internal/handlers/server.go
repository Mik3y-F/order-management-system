package handlers

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	pb "github.com/Mik3y-F/order-management-system/orders/api/generated"
	"github.com/Mik3y-F/order-management-system/orders/internal/service"
	"google.golang.org/grpc"
)

// GRPCServer struct represents the GRPC server for the order management system.
type GRPCServer struct {
	pb.UnimplementedOrdersServer

	grpcServer *grpc.Server
	mu         sync.Mutex // synchronizes access to the grpcServer

	// Internal services
	ProductService  service.ProductService
	CustomerService service.CustomerService
	OrderService    service.OrderService
}

// NewGRPCServer creates a new instance of GRPCServer.
func NewGRPCServer() *GRPCServer {
	return &GRPCServer{}
}

// Run starts the GRPC server on the specified bind address and port.
func (s *GRPCServer) Run(ctx context.Context, bindAddress string, port string) error {
	lis, err := net.Listen("tcp", bindAddress+":"+port)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	s.mu.Lock()
	s.grpcServer = grpc.NewServer()
	s.mu.Unlock()

	pb.RegisterOrdersServer(s.grpcServer, s)

	log.Printf("Starting server on port %v", lis.Addr())

	return s.grpcServer.Serve(lis)
}

// Stop gracefully stops the GRPC server.
func (s *GRPCServer) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.grpcServer != nil {
		s.grpcServer.GracefulStop()
	}
}
