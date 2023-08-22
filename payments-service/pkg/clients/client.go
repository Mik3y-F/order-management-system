package payments

import (
	"context"
	"log"

	pb "github.com/Mik3y-F/order-management-system/payments/api/generated"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type HealthCheckRequest = pb.HealthCheckRequest
type HealthCheckResponse = pb.HealthCheckResponse

type PaymentsClient interface {
	HealthCheck(ctx context.Context, product *HealthCheckRequest) (*HealthCheckResponse, error)
}

type GrpcPaymentsClient struct {
	conn   *grpc.ClientConn
	client pb.PaymentsClient
}

func NewGrpcPaymentsClient(conn *grpc.ClientConn) *GrpcPaymentsClient {

	client := pb.NewPaymentsClient(conn)

	return &GrpcPaymentsClient{
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
