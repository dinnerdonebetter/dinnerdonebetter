//go:build wireinject

package searchdataindexscheduler

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/config"
	"github.com/dinnerdonebetter/backend/internal/repositories/postgres/auditlogentries"
	"github.com/dinnerdonebetter/backend/internal/repositories/postgres/identity"
	"github.com/dinnerdonebetter/backend/internal/repositories/postgres/mealplanning"

	"github.com/google/wire"
	databasecfg "github.com/verygoodsoftwarenotvirus/platform/database/config"
	"github.com/verygoodsoftwarenotvirus/platform/database/postgres"
	msgconfig "github.com/verygoodsoftwarenotvirus/platform/messagequeue/config"
	"github.com/verygoodsoftwarenotvirus/platform/observability"
	loggingcfg "github.com/verygoodsoftwarenotvirus/platform/observability/logging/config"
	metricscfg "github.com/verygoodsoftwarenotvirus/platform/observability/metrics/config"
	tracingcfg "github.com/verygoodsoftwarenotvirus/platform/observability/tracing/config"
	"github.com/verygoodsoftwarenotvirus/platform/search/text/indexing"
)

// Build builds a server.
func Build(
	ctx context.Context,
	cfg *config.SearchDataIndexSchedulerConfig,
) (*indexing.IndexScheduler, error) {
	wire.Build(
		indexing.ProvidersIndexing,
		tracingcfg.TracingConfigProviders,
		observability.O11yProviders,
		msgconfig.MessageQueueProviders,
		databasecfg.ClientConfigProviders,
		postgres.PGProviders,
		loggingcfg.LogConfigProviders,
		metricscfg.MetricsConfigProviders,
		auditlogentries.AuditRepoProviders,
		identity.IDRepoProviders,
		mealplanning.MPRepoProviders,
		ProvideIndexFunctions,
		ConfigProviders,
	)

	return nil, nil
}
