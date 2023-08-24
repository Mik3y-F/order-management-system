package client

import (
	pb "github.com/Mik3y-F/order-management-system/orders/api/generated"
)

type HealthCheckRequest = pb.HealthCheckRequest
type HealthCheckResponse = pb.HealthCheckResponse

type UpdateOrderStatusRequest = pb.UpdateOrderStatusRequest
type UpdateOrderStatusResponse = pb.UpdateOrderStatusResponse

var OrderStatusPaid = pb.OrderStatus_PAID
var OrderStatusCancelled = pb.OrderStatus_CANCELLED
var OrderStatusFailed = pb.OrderStatus_FAILED
var OrderStatusPending = pb.OrderStatus_PENDING
