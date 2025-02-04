//go:build wireinject
// +build wireinject

package searchdataindexscheduler

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/config"
	"github.com/dinnerdonebetter/backend/internal/database/postgres"
	msgconfig "github.com/dinnerdonebetter/backend/internal/lib/messagequeue/config"
	"github.com/dinnerdonebetter/backend/internal/lib/observability"
	loggingcfg "github.com/dinnerdonebetter/backend/internal/lib/observability/logging/config"
	metricscfg "github.com/dinnerdonebetter/backend/internal/lib/observability/metrics/config"
	tracingcfg "github.com/dinnerdonebetter/backend/internal/lib/observability/tracing/config"
	"github.com/dinnerdonebetter/backend/internal/lib/search/text/indexing"

	"github.com/google/wire"
)

// Build builds a server.
func Build(
	ctx context.Context,
	cfg *config.SearchDataIndexSchedulerConfig,
) (*indexing.IndexScheduler, error) {
	wire.Build(
		indexing.ProvidersIndexing,
		tracingcfg.ProvidersTracingConfig,
		observability.ProvidersObservability,
		msgconfig.MessageQueueProviders,
		postgres.ProvidersPostgres,
		loggingcfg.ProvidersLoggingConfig,
		metricscfg.ProvidersMetrics,
		ProvideIndexFunctions,
		ConfigProviders,
	)

	return nil, nil
}
