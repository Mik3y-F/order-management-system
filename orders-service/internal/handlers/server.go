package handlers

import (
	"log"
	"net"

	pb "github.com/Mik3y-F/order-management-system/orders/api/generated"

	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

type GRPCServer struct {
	pb.UnimplementedOrdersServer
}

func NewGRPCServer() *GRPCServer {
	return &GRPCServer{}
}

func (s *GRPCServer) Run(bindAddress string, port string) error {
	lis, err := net.Listen("tcp", bindAddress+":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterOrdersServer(grpcServer, s)

	log.Printf("Starting server on port %v", lis.Addr())

	return grpcServer.Serve(lis)
}
