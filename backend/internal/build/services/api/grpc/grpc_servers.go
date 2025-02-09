package grpcapi

import (
	"github.com/dinnerdonebetter/backend/internal/grpc/service"
	"github.com/dinnerdonebetter/backend/internal/grpcimpl/serverimpl"
	"github.com/dinnerdonebetter/backend/internal/lib/server/grpc"

	grpc2 "google.golang.org/grpc"
)

func BuildRegistrationFuncs(eatingServer *serverimpl.Server) []grpc.RegistrationFunc {
	return []grpc.RegistrationFunc{
		func(server *grpc2.Server) {
			service.RegisterEatingServiceServer(server, eatingServer)
		},
	}
}

func BuildUnaryServerInterceptors(eatingServer *serverimpl.Server) []grpc2.UnaryServerInterceptor {
	return []grpc2.UnaryServerInterceptor{
		eatingServer.AuthInterceptor(),
	}
}

func BuildStreamServerInterceptors(eatingServer *serverimpl.Server) []grpc2.StreamServerInterceptor {
	return []grpc2.StreamServerInterceptor{
		//
	}
}
