package grpc

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	pb "github.com/Mik3y-F/order-management-system/payments/api/generated"
	"github.com/Mik3y-F/order-management-system/payments/internal/service"

	"google.golang.org/grpc"
)

type GRPCServer struct {
	pb.UnimplementedPaymentsServer

	grpcServer *grpc.Server
	mu         sync.Mutex // synchronizes access to the grpcServer

	// Internal servicesx
	PaymentsService service.PaymentsService
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

	pb.RegisterPaymentsServer(s.grpcServer, s)

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
