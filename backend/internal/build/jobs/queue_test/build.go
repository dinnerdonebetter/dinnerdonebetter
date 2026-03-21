package queuetest

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/config"
	"github.com/dinnerdonebetter/backend/internal/repositories/postgres/internalops"
	queuetest "github.com/dinnerdonebetter/backend/internal/services/internalops/workers/queue_test"

	"github.com/samber/do/v2"
	"github.com/verygoodsoftwarenotvirus/platform/database/postgres"
	msgconfig "github.com/verygoodsoftwarenotvirus/platform/messagequeue/config"
	"github.com/verygoodsoftwarenotvirus/platform/observability"
	loggingcfg "github.com/verygoodsoftwarenotvirus/platform/observability/logging/config"
	"github.com/verygoodsoftwarenotvirus/platform/observability/metrics"
	metricscfg "github.com/verygoodsoftwarenotvirus/platform/observability/metrics/config"
	tracingcfg "github.com/verygoodsoftwarenotvirus/platform/observability/tracing/config"
)

// Build builds the queue test job and a cleanup that flushes metrics.
func Build(
	ctx context.Context,
	cfg *config.QueueTestJobConfig,
) (*BuildResult, error) {
	i := do.New()

	do.ProvideValue(i, ctx)
	do.ProvideValue(i, cfg)

	RegisterConfigs(i)

	observability.RegisterO11yConfigs(i)
	tracingcfg.RegisterTracerProvider(i)
	loggingcfg.RegisterLogger(i)
	metricscfg.RegisterMetricsProvider(i)
	msgconfig.RegisterMessageQueue(i)
	postgres.RegisterDatabaseClient(i)
	internalops.RegisterInternalOpsRepository(i)
	queuetest.RegisterQueueTest(i)

	do.Provide[*BuildResult](i, func(i do.Injector) (*BuildResult, error) {
		return NewBuildResult(
			do.MustInvoke[*queuetest.Job](i),
			do.MustInvoke[metrics.Provider](i),
		), nil
	})

	return do.MustInvoke[*BuildResult](i), nil
}
