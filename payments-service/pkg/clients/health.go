package payments

import (
	"context"
)

func (c *GrpcPaymentsClient) CreateProduct(ctx context.Context, req *HealthCheckRequest) (*HealthCheckResponse, error) {
	return c.client.HealthCheck(ctx, req)
}
