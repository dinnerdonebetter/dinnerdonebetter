package grpcapi

import (
	"github.com/dinnerdonebetter/backend/internal/grpc/service"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/lib/server/grpc"

	grpc2 "google.golang.org/grpc"
)

func BuildWrappedServer(cfg *grpc.Config, logger logging.Logger, eatingServer service.EatingServiceServer) (*grpc.Server, error) {
	server, err := grpc.NewGRPCServer(cfg, logger, func(server *grpc2.Server) {
		service.RegisterEatingServiceServer(server, eatingServer)
	})
	if err != nil {
		return nil, err
	}

	return server, nil
}

func BuildRegistrationFuncs(eatingServer service.EatingServiceServer) []grpc.RegistrationFunc {
	return []grpc.RegistrationFunc{
		func(server *grpc2.Server) {
			service.RegisterEatingServiceServer(server, eatingServer)
		},
	}
}
