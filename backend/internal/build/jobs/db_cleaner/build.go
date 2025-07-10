//go:build wireinject
// +build wireinject

package dbcleaner

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/config"
	"github.com/dinnerdonebetter/backend/internal/platform/database/postgres"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	loggingcfg "github.com/dinnerdonebetter/backend/internal/platform/observability/logging/config"
	metricscfg "github.com/dinnerdonebetter/backend/internal/platform/observability/metrics/config"
	tracingcfg "github.com/dinnerdonebetter/backend/internal/platform/observability/tracing/config"
	"github.com/dinnerdonebetter/backend/internal/repositories/postgres/maintenance"
	dbcleaner "github.com/dinnerdonebetter/backend/internal/services/core/workers/db_cleaner"

	"github.com/google/wire"
)

// Build builds a server.
func Build(
	ctx context.Context,
	cfg *config.DBCleanerConfig,
) (*dbcleaner.Job, error) {
	wire.Build(
		dbcleaner.ProvidersDBCleaner,
		tracingcfg.ProvidersTracingConfig,
		observability.Providers,
		postgres.Providers,
		loggingcfg.ProvidersLogConfig,
		metricscfg.Providers,
		maintenance.Providers,
		ConfigProviders,
	)

	return nil, nil
}
