package grpcapi

import (
	auditsvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/audit"
	"github.com/dinnerdonebetter/backend/internal/platform/server/grpc"

	grpc2 "google.golang.org/grpc"
)

func BuildRegistrationFuncs(auditLogService auditsvc.AuditServiceServer) []grpc.RegistrationFunc {
	return []grpc.RegistrationFunc{
		func(server *grpc2.Server) {
			auditsvc.RegisterAuditServiceServer(server, auditLogService)
		},
	}
}

func BuildUnaryServerInterceptors() []grpc2.UnaryServerInterceptor {
	return []grpc2.UnaryServerInterceptor{
		//eatingServer.AuthInterceptor(),
	}
}

func BuildStreamServerInterceptors() []grpc2.StreamServerInterceptor {
	return []grpc2.StreamServerInterceptor{
		//
	}
}
