package client

import (
	pb "github.com/Mik3y-F/order-management-system/payments/api/generated"
)

type HealthCheckRequest = pb.HealthCheckRequest
type HealthCheckResponse = pb.HealthCheckResponse

type ProcessMpesaPaymentRequest = pb.MpesaPaymentRequest
type ProcessMpesaPaymentResponse = pb.MpesaPaymentResponse
