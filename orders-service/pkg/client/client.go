package client

import (
	"context"
	"log"

	pb "github.com/Mik3y-F/order-management-system/orders/api/generated"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var _ OrdersClient = (*GrpcOrderClient)(nil)

type OrdersClient interface {
	HealthCheck(ctx context.Context, product *HealthCheckRequest) (*HealthCheckResponse, error)
	UpdateOrderStatus(ctx context.Context, req *pb.UpdateOrderStatusRequest) (*pb.UpdateOrderStatusResponse, error)
}
type GrpcOrderClient struct {
	conn   *grpc.ClientConn
	client pb.OrdersClient
}

func NewGrpcOrderClient(conn *grpc.ClientConn) *GrpcOrderClient {

	client := pb.NewOrdersClient(conn)

	return &GrpcOrderClient{
		conn:   conn,
		client: client,
	}
}

func ConnectToOrderService(address string) (*grpc.ClientConn, error) {

	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	return conn, nil
}
