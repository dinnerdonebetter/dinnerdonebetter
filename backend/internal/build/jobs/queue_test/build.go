package queuetest

import (
	"context"

	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/config"
	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/repositories/postgres/internalops"
	queuetest "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/services/internalops/workers/queue_test"

	"github.com/samber/do/v2"
	databasecfg "github.com/verygoodsoftwarenotvirus/platform/v3/database/config"
	"github.com/verygoodsoftwarenotvirus/platform/v3/database/postgres"
	msgconfig "github.com/verygoodsoftwarenotvirus/platform/v3/messagequeue/config"
	"github.com/verygoodsoftwarenotvirus/platform/v3/observability"
	loggingcfg "github.com/verygoodsoftwarenotvirus/platform/v3/observability/logging/config"
	"github.com/verygoodsoftwarenotvirus/platform/v3/observability/metrics"
	metricscfg "github.com/verygoodsoftwarenotvirus/platform/v3/observability/metrics/config"
	tracingcfg "github.com/verygoodsoftwarenotvirus/platform/v3/observability/tracing/config"
)

// BuildInjector creates and configures the dependency injection container.
func BuildInjector(
	ctx context.Context,
	cfg *config.QueueTestJobConfig,
) *do.RootScope {
	i := do.New()

	do.ProvideValue(i, ctx)
	do.ProvideValue(i, cfg)

	RegisterConfigs(i)

	observability.RegisterO11yConfigs(i)
	tracingcfg.RegisterTracerProvider(i)
	loggingcfg.RegisterLogger(i)
	metricscfg.RegisterMetricsProvider(i)
	msgconfig.RegisterMessageQueue(i)
	databasecfg.RegisterClientConfig(i)
	postgres.RegisterDatabaseClient(i)
	internalops.RegisterInternalOpsRepository(i)
	queuetest.RegisterQueueTest(i)

	do.Provide[*BuildResult](i, func(i do.Injector) (*BuildResult, error) {
		return NewBuildResult(
			do.MustInvoke[*queuetest.Job](i),
			do.MustInvoke[metrics.Provider](i),
		), nil
	})

	return i
}

// Build builds the queue test job and a cleanup that flushes metrics.
func Build(
	ctx context.Context,
	cfg *config.QueueTestJobConfig,
) (*BuildResult, error) {
	i := BuildInjector(ctx, cfg)
	return do.MustInvoke[*BuildResult](i), nil
}
