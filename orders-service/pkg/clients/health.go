package orders

import (
	"context"
)

func (c *GrpcOrderClient) CreateProduct(ctx context.Context, req *HealthCheckRequest) (*HealthCheckResponse, error) {
	return c.client.HealthCheck(ctx, req)
}
