//go:build wireinject
// +build wireinject

package grpcapi

import (
	"context"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"

	"github.com/dinnerdonebetter/backend/internal/config"
	"github.com/dinnerdonebetter/backend/internal/platform/database/postgres"
	loggingcfg "github.com/dinnerdonebetter/backend/internal/platform/observability/logging/config"
	tracingcfg "github.com/dinnerdonebetter/backend/internal/platform/observability/tracing/config"
	auditrepo "github.com/dinnerdonebetter/backend/internal/repositories/postgres/auditlogentries"
	auditsvc "github.com/dinnerdonebetter/backend/internal/services/audit/grpc"

	"github.com/google/wire"
)

// Build builds a server.
func Build(
	ctx context.Context,
	cfg *config.APIServiceConfig,
) (*GRPCService, error) {
	wire.Build(ConfigProviders,
		loggingcfg.ProvidersLogConfig,
		tracingcfg.ProvidersTracingConfig,
		observability.Providers,
		postgres.Providers,
		auditsvc.NewService,
		auditrepo.Providers,
		// BuildUnaryServerInterceptors,
		// BuildStreamServerInterceptors,
		// BuildRegistrationFuncs,
		NewGRPCService,
	)

	return nil, nil
}
