//go:build wireinject
// +build wireinject

package api

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/config"
	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/database/postgres"
	analyticscfg "github.com/dinnerdonebetter/backend/internal/lib/analytics/config"
	"github.com/dinnerdonebetter/backend/internal/lib/authentication"
	"github.com/dinnerdonebetter/backend/internal/lib/encoding"
	featureflagscfg "github.com/dinnerdonebetter/backend/internal/lib/featureflags/config"
	msgconfig "github.com/dinnerdonebetter/backend/internal/lib/messagequeue/config"
	"github.com/dinnerdonebetter/backend/internal/lib/observability"
	loggingcfg "github.com/dinnerdonebetter/backend/internal/lib/observability/logging/config"
	metricscfg "github.com/dinnerdonebetter/backend/internal/lib/observability/metrics/config"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
	tracingcfg "github.com/dinnerdonebetter/backend/internal/lib/observability/tracing/config"
	routingcfg "github.com/dinnerdonebetter/backend/internal/lib/routing/config"
	"github.com/dinnerdonebetter/backend/internal/lib/server/http"
	authservice "github.com/dinnerdonebetter/backend/internal/services/core/handlers/authentication"

	"github.com/google/wire"
)

// Build builds a server.
func Build(
	ctx context.Context,
	cfg *config.APIServiceConfig,
) (http.Server, error) {
	wire.Build(
		authentication.AuthProviders,
		database.DBProviders,
		encoding.EncDecProviders,
		msgconfig.MessageQueueProviders,
		analyticscfg.ProvidersAnalytics,
		featureflagscfg.ProvidersFeatureFlags,
		tracing.ProvidersTracing,
		tracingcfg.ProvidersTracingConfig,
		observability.ProvidersObservability,
		postgres.ProvidersPostgres,
		loggingcfg.ProvidersLoggingConfig,
		authservice.Providers,
		metricscfg.ProvidersMetrics,
		http.ProvidersHTTP,
		routingcfg.RoutingConfigProviders,
		ConfigProviders,
		ProvideAPIRouter,
	)

	return nil, nil
}
