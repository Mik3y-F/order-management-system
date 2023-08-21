package handlers_test

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/Mik3y-F/order-management-system/orders/internal/handlers"
	grpc_handlers "github.com/Mik3y-F/order-management-system/orders/internal/handlers"
	"github.com/Mik3y-F/order-management-system/orders/internal/mock"
)

const (
	INVALID_PORT = "70000"
)

type TestGRPCServer struct {
	*grpc_handlers.GRPCServer

	// Add mock services here
	ProductService  mock.ProductService
	CustomerService mock.CustomerService
}

func NewTestGRPCServer(tb testing.TB) *TestGRPCServer {
	s := &TestGRPCServer{
		GRPCServer: grpc_handlers.NewGRPCServer(),
	}

	// Set mock services here
	s.GRPCServer.ProductService = &s.ProductService
	s.GRPCServer.CustomerService = &s.CustomerService

	return s
}

func TestGRPCServer_Run(t *testing.T) {
	tests := []struct {
		name     string
		bindAddr string
		port     string
		wantErr  bool
	}{
		{
			name:     "valid bind address and port",
			bindAddr: "localhost",
			port:     getFreePort(t),
			wantErr:  false,
		},
		{
			name:     "invalid port",
			bindAddr: "localhost",
			port:     INVALID_PORT,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := handlers.NewGRPCServer()
			errCh := make(chan error)

			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			go func() {
				errCh <- server.Run(ctx, tt.bindAddr, tt.port)
			}()

			// Wait for server to respond or timeout after a specified time
			select {
			case err := <-errCh:
				if (err != nil) != tt.wantErr {
					t.Errorf("GRPCServer.Run() error = %v, wantErr %v", err, tt.wantErr)
				}
			case <-ctx.Done():
				if tt.wantErr {
					t.Errorf("GRPCServer.Run() expected error for port %v, but got none", tt.port)
				}
			}

			server.Stop()
		})
	}
}

// Helper function to get a free port
func getFreePort(t *testing.T) string {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("error resolving tcp address: %v", err)
	}

	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		t.Fatalf("error listening on tcp address: %v", err)
	}

	defer listener.Close()
	return fmt.Sprintf("%d", listener.Addr().(*net.TCPAddr).Port)
}
