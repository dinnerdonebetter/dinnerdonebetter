//go:build wireinject

package queuetest

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/config"
	databasecfg "github.com/dinnerdonebetter/backend/internal/platform/database/config"
	"github.com/dinnerdonebetter/backend/internal/platform/database/postgres"
	msgconfig "github.com/dinnerdonebetter/backend/internal/platform/messagequeue/config"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	loggingcfg "github.com/dinnerdonebetter/backend/internal/platform/observability/logging/config"
	metricscfg "github.com/dinnerdonebetter/backend/internal/platform/observability/metrics/config"
	tracingcfg "github.com/dinnerdonebetter/backend/internal/platform/observability/tracing/config"
	"github.com/dinnerdonebetter/backend/internal/repositories/postgres/internalops"
	queuetest "github.com/dinnerdonebetter/backend/internal/services/internalops/workers/queue_test"

	"github.com/google/wire"
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
