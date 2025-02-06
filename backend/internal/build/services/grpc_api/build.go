//go:build wireinject
// +build wireinject

package grpcapi

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/config"
	"github.com/dinnerdonebetter/backend/internal/database/postgres"
	"github.com/dinnerdonebetter/backend/internal/grpcimpl/serverimpl"
	analyticscfg "github.com/dinnerdonebetter/backend/internal/lib/analytics/config"
	"github.com/dinnerdonebetter/backend/internal/lib/authentication"
	tokenscfg "github.com/dinnerdonebetter/backend/internal/lib/authentication/tokens/config"
	featureflagscfg "github.com/dinnerdonebetter/backend/internal/lib/featureflags/config"
	msgconfig "github.com/dinnerdonebetter/backend/internal/lib/messagequeue/config"
	"github.com/dinnerdonebetter/backend/internal/lib/observability"
	loggingcfg "github.com/dinnerdonebetter/backend/internal/lib/observability/logging/config"
	metricscfg "github.com/dinnerdonebetter/backend/internal/lib/observability/metrics/config"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
	tracingcfg "github.com/dinnerdonebetter/backend/internal/lib/observability/tracing/config"
	"github.com/dinnerdonebetter/backend/internal/lib/random"
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
		authentication.AuthProviders,
		grpc.ProvidersGRPC,
		random.ProvidersRandom,
		msgconfig.MessageQueueProviders,
		analyticscfg.ProvidersAnalytics,
		tokenscfg.ProvidersTokenIssuers,
		metricscfg.ProvidersMetrics,
		featureflagscfg.ProvidersFeatureFlags,
		tracing.ProvidersTracing,
		tracingcfg.ProvidersTracingConfig,
		observability.ProvidersObservability,
		postgres.ProvidersPostgres,
		ConfigProviders,
		BuildUnaryServerInterceptors,
		BuildStreamServerInterceptors,
		BuildRegistrationFuncs,
	)

	return nil, nil
}
