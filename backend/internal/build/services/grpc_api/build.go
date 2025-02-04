//go:build wireinject
// +build wireinject

package grpcapi

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/config"
	"github.com/dinnerdonebetter/backend/internal/database/postgres"
	"github.com/dinnerdonebetter/backend/internal/grpcimpl/serverimpl"
	"github.com/dinnerdonebetter/backend/internal/lib/observability"
	loggingcfg "github.com/dinnerdonebetter/backend/internal/lib/observability/logging/config"
	tracingcfg "github.com/dinnerdonebetter/backend/internal/lib/observability/tracing/config"
	"github.com/dinnerdonebetter/backend/internal/lib/server/grpc"

	"github.com/google/wire"
)

// Build builds a server.
func Build(
	ctx context.Context,
	cfg *config.APIServiceConfig,
) (*grpc.Server, error) {
	wire.Build(
		serverimpl.Providers,
		loggingcfg.ProvidersLoggingConfig,
		grpc.ProvidersGRPC,
		ConfigProviders,
		BuildUnaryServerInterceptors,
		BuildStreamServerInterceptors,
		BuildRegistrationFuncs,
		tracingcfg.ProvidersTracingConfig,
		observability.ProvidersObservability,
		postgres.ProvidersPostgres,
	)

	return nil, nil
}
