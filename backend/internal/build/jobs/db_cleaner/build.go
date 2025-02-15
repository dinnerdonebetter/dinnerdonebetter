//go:build wireinject
// +build wireinject

package dbcleaner

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/config"
	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/database/postgres"
	"github.com/dinnerdonebetter/backend/internal/lib/observability"
	loggingcfg "github.com/dinnerdonebetter/backend/internal/lib/observability/logging/config"
	metricscfg "github.com/dinnerdonebetter/backend/internal/lib/observability/metrics/config"
	tracingcfg "github.com/dinnerdonebetter/backend/internal/lib/observability/tracing/config"
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
		database.DBProviders,
		tracingcfg.ProvidersTracingConfig,
		observability.ProvidersObservability,
		postgres.ProvidersPostgres,
		loggingcfg.ProvidersLogConfig,
		metricscfg.ProvidersMetrics,
		ConfigProviders,
	)

	return nil, nil
}
