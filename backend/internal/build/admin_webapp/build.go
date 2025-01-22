//go:build wireinject
// +build wireinject

package adminwebapp

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/config"
	"github.com/dinnerdonebetter/backend/internal/observability"
	loggingcfg "github.com/dinnerdonebetter/backend/internal/observability/logging/config"
	metricscfg "github.com/dinnerdonebetter/backend/internal/observability/metrics/config"
	tracingcfg "github.com/dinnerdonebetter/backend/internal/observability/tracing/config"
	"github.com/dinnerdonebetter/backend/internal/server/http"
	"github.com/dinnerdonebetter/backend/internal/services/frontend/admin"

	"github.com/google/wire"
)

// Build builds a server.
func Build(
	ctx context.Context,
	cfg *config.AdminWebappConfig,
) (http.Server, error) {
	wire.Build(
		admin.ProvidersAdminWebapp,
		tracingcfg.ProvidersTracingConfig,
		observability.ProvidersObservability,
		loggingcfg.ProvidersLogConfig,
		metricscfg.ProvidersMetrics,
		http.ProvidersHTTP,
		ConfigProviders,
		ProvideAdminWebappRouter,
	)

	return nil, nil
}
