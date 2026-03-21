//go:build wireinject

package dbcleaner

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/config"
	"github.com/dinnerdonebetter/backend/internal/repositories/postgres/internalops"
	dbcleaner "github.com/dinnerdonebetter/backend/internal/services/oauth/workers/db_cleaner"

	"github.com/google/wire"
	databasecfg "github.com/verygoodsoftwarenotvirus/platform/database/config"
	"github.com/verygoodsoftwarenotvirus/platform/database/postgres"
	"github.com/verygoodsoftwarenotvirus/platform/observability"
	loggingcfg "github.com/verygoodsoftwarenotvirus/platform/observability/logging/config"
	metricscfg "github.com/verygoodsoftwarenotvirus/platform/observability/metrics/config"
	tracingcfg "github.com/verygoodsoftwarenotvirus/platform/observability/tracing/config"
)

// Build builds a server.
func Build(
	ctx context.Context,
	cfg *config.DBCleanerConfig,
) (*dbcleaner.Job, error) {
	wire.Build(
		dbcleaner.ProvidersDBCleaner,
		tracingcfg.TracingConfigProviders,
		observability.O11yProviders,
		databasecfg.ClientConfigProviders,
		postgres.PGProviders,
		loggingcfg.LogConfigProviders,
		metricscfg.MetricsConfigProviders,
		internalops.Providers,
		ConfigProviders,
	)

	return nil, nil
}
