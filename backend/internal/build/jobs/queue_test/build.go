//go:build wireinject

package queuetest

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/config"
	"github.com/dinnerdonebetter/backend/internal/repositories/postgres/internalops"
	queuetest "github.com/dinnerdonebetter/backend/internal/services/internalops/workers/queue_test"

	"github.com/google/wire"
	databasecfg "github.com/verygoodsoftwarenotvirus/platform/database/config"
	"github.com/verygoodsoftwarenotvirus/platform/database/postgres"
	msgconfig "github.com/verygoodsoftwarenotvirus/platform/messagequeue/config"
	"github.com/verygoodsoftwarenotvirus/platform/observability"
	loggingcfg "github.com/verygoodsoftwarenotvirus/platform/observability/logging/config"
	metricscfg "github.com/verygoodsoftwarenotvirus/platform/observability/metrics/config"
	tracingcfg "github.com/verygoodsoftwarenotvirus/platform/observability/tracing/config"
)

// Build builds the queue test job and a cleanup that flushes metrics.
func Build(
	ctx context.Context,
	cfg *config.QueueTestJobConfig,
) (*BuildResult, error) {
	wire.Build(
		NewBuildResult,
		queuetest.ProvidersQueueTest,
		tracingcfg.TracingConfigProviders,
		observability.O11yProviders,
		loggingcfg.LogConfigProviders,
		metricscfg.MetricsConfigProviders,
		msgconfig.MessageQueueProviders,
		databasecfg.ClientConfigProviders,
		postgres.PGProviders,
		internalops.Providers,
		ConfigProviders,
	)

	return nil, nil
}
