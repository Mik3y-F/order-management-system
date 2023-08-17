package handlers

import (
	"context"
	"reflect"
	"testing"

	pb "github.com/Mik3y-F/order-management-system/orders/api/generated"
)

func TestGRPCServer_HealthCheck(t *testing.T) {

	s := NewGRPCServer()

	type args struct {
		ctx context.Context
		in  *pb.HealthCheckRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *pb.HealthCheckResponse
		wantErr bool
	}{
		{
			name: "Test HealthCheck",
			args: args{
				ctx: context.Background(),
				in:  &pb.HealthCheckRequest{},
			},
			want: &pb.HealthCheckResponse{
				Status: "OK",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := s.HealthCheck(tt.args.ctx, tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("GRPCServer.HealthCheck() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GRPCServer.HealthCheck() = %v, want %v", got, tt.want)
			}
		})
	}
}
