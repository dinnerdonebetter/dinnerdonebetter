package grpcapi

import (
	"context"
	auditsvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/audit"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/metrics"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	textsearchcfg "github.com/dinnerdonebetter/backend/internal/platform/search/text/config"
	"github.com/dinnerdonebetter/backend/internal/platform/server/grpc"
	identityindexing "github.com/dinnerdonebetter/backend/internal/services/identity/indexing"

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

func ProvideUserTextSearcher(
	ctx context.Context,
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	metricsProvider metrics.Provider,
	cfg *textsearchcfg.Config,
) (identityindexing.UserTextSearcher, error) {
	return textsearchcfg.ProvideIndex[identityindexing.UserSearchSubset](
		ctx,
		logger,
		tracerProvider, metricsProvider,
		cfg,
		identityindexing.IndexTypeUsers,
	)
}
