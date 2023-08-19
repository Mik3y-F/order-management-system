package handlers

import (
	"fmt"
	"log"
	"net"

	pb "github.com/Mik3y-F/order-management-system/orders/api/generated"

	"github.com/Mik3y-F/order-management-system/orders/internal/service"

	"google.golang.org/grpc"
)

type GRPCServer struct {
	pb.UnimplementedOrdersServer

	ProductService service.ProductService
}

func NewGRPCServer() *GRPCServer {
	return &GRPCServer{}
}

func (s *GRPCServer) Run(bindAddress string, port string) error {
	lis, err := net.Listen("tcp", bindAddress+":"+port)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterOrdersServer(grpcServer, s)

	log.Printf("Starting server on port %v", lis.Addr())

	return grpcServer.Serve(lis)
}
