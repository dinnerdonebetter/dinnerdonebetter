package dbcleaner

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/config"
	"github.com/dinnerdonebetter/backend/internal/repositories/postgres/internalops"
	dbcleaner "github.com/dinnerdonebetter/backend/internal/services/oauth/workers/db_cleaner"

	"github.com/samber/do/v2"
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
	i := do.New()

	do.ProvideValue(i, ctx)
	do.ProvideValue(i, cfg)

	RegisterConfigs(i)

	observability.RegisterO11yConfigs(i)
	tracingcfg.RegisterTracerProvider(i)
	loggingcfg.RegisterLogger(i)
	metricscfg.RegisterMetricsProvider(i)
	postgres.RegisterDatabaseClient(i)
	internalops.RegisterInternalOpsRepository(i)
	dbcleaner.RegisterDBCleaner(i)

	return do.MustInvoke[*dbcleaner.Job](i), nil
}
