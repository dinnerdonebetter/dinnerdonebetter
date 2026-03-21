package searchdataindexscheduler

import (
	"context"

	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/config"
	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/identity"
	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/repositories/postgres/auditlogentries"
	identityrepo "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/repositories/postgres/identity"
	mealplanningrepo "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/repositories/postgres/mealplanning"

	"github.com/samber/do/v2"
	databasecfg "github.com/verygoodsoftwarenotvirus/platform/database/config"
	"github.com/verygoodsoftwarenotvirus/platform/database/postgres"
	msgconfig "github.com/verygoodsoftwarenotvirus/platform/messagequeue/config"
	"github.com/verygoodsoftwarenotvirus/platform/observability"
	loggingcfg "github.com/verygoodsoftwarenotvirus/platform/observability/logging/config"
	metricscfg "github.com/verygoodsoftwarenotvirus/platform/observability/metrics/config"
	tracingcfg "github.com/verygoodsoftwarenotvirus/platform/observability/tracing/config"
	"github.com/verygoodsoftwarenotvirus/platform/search/text/indexing"
)

// BuildInjector creates and configures the dependency injection container.
func BuildInjector(
	ctx context.Context,
	cfg *config.SearchDataIndexSchedulerConfig,
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
	auditlogentries.RegisterAuditLogRepository(i)
	identityrepo.RegisterIdentityRepository(i)
	mealplanningrepo.RegisterMealPlanningRepository(i)

	do.Provide[map[string]indexing.Function](i, func(i do.Injector) (map[string]indexing.Function, error) {
		identityRepo := do.MustInvoke[identity.Repository](i)
		mealPlanningRepo := do.MustInvoke[mealplanning.Repository](i)
		return ProvideIndexFunctions(identityRepo, mealPlanningRepo), nil
	})

	indexing.RegisterIndexScheduler(i)

	return i
}

// Build builds a server.
func Build(
	ctx context.Context,
	cfg *config.SearchDataIndexSchedulerConfig,
) (*indexing.IndexScheduler, error) {
	i := BuildInjector(ctx, cfg)
	return do.MustInvoke[*indexing.IndexScheduler](i), nil
}
