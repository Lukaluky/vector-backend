package grpc

import (
	"context"
	"fmt"
	"net"

	gogrpc "google.golang.org/grpc"

	shipmentpb "vektor-backend/proto"
)

type Server struct {
	port   string
	server *gogrpc.Server
}

func NewServer(port string, handler *Handler) *Server {
	s := gogrpc.NewServer()

	shipmentpb.RegisterShipmentServiceServer(s, handler)

	return &Server{
		port:   port,
		server: s,
	}
}

func (s *Server) Run() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", s.port))
	if err != nil {
		return fmt.Errorf("listen tcp: %w", err)
	}

	if err := s.server.Serve(lis); err != nil {
		return fmt.Errorf("serve grpc: %w", err)
	}

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	done := make(chan struct{})

	go func() {
		s.server.GracefulStop()
		close(done)
	}()

	select {
	case <-ctx.Done():
		s.server.Stop()
		return ctx.Err()
	case <-done:
		return nil
	}
}