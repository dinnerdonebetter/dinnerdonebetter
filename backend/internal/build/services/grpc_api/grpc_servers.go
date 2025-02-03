package grpcapi

import (
	"github.com/dinnerdonebetter/backend/internal/grpc/service"
	"github.com/dinnerdonebetter/backend/internal/lib/server/grpc"

	grpc2 "google.golang.org/grpc"
)

func BuildWrappedServer(cfg *grpc.Config, eatingServer service.EatingServiceServer) (*grpc.Server, error) {
	server, err := grpc.NewGRPCServer(cfg, func(server *grpc2.Server) {
		service.RegisterEatingServiceServer(server, eatingServer)
	})
	if err != nil {
		return nil, err
	}

	return server, nil
}
