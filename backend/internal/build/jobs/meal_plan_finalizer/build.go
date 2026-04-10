package mealplanfinalizer

import (
	"context"

	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/config"
	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/repositories/postgres/auditlogentries"
	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/repositories/postgres/identity"
	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/repositories/postgres/mealplanning"
	mealplanfinalizer "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/services/mealplanning/workers/meal_plan_finalizer"

	databasecfg "github.com/verygoodsoftwarenotvirus/platform/v5/database/config"
	"github.com/verygoodsoftwarenotvirus/platform/v5/database/postgres"
	msgconfig "github.com/verygoodsoftwarenotvirus/platform/v5/messagequeue/config"
	"github.com/verygoodsoftwarenotvirus/platform/v5/observability"
	loggingcfg "github.com/verygoodsoftwarenotvirus/platform/v5/observability/logging/config"
	metricscfg "github.com/verygoodsoftwarenotvirus/platform/v5/observability/metrics/config"
	tracingcfg "github.com/verygoodsoftwarenotvirus/platform/v5/observability/tracing/config"

	"github.com/samber/do/v2"
)

// BuildInjector creates and configures the dependency injection container.
func BuildInjector(
	ctx context.Context,
	cfg *config.MealPlanFinalizerConfig,
) *do.RootScope {
	i := do.New()

	do.ProvideValue(i, ctx)
	do.ProvideValue(i, cfg)

	RegisterConfigs(i)

	observability.RegisterO11yConfigs(i)
	tracingcfg.RegisterTracerProvider(i)
	loggingcfg.RegisterLogger(i)
	metricscfg.RegisterMetricsProvider(i)
	databasecfg.RegisterClientConfig(i)
	postgres.RegisterDatabaseClient(i)
	msgconfig.RegisterMessageQueue(i)
	auditlogentries.RegisterAuditLogRepository(i)
	identity.RegisterIdentityRepository(i)
	mealplanning.RegisterMealPlanningRepository(i)
	mealplanfinalizer.RegisterMealPlanFinalizer(i)

	return i
}

// Build builds a server.
func Build(
	ctx context.Context,
	cfg *config.MealPlanFinalizerConfig,
) (*mealplanfinalizer.Worker, error) {
	i := BuildInjector(ctx, cfg)
	return do.MustInvoke[*mealplanfinalizer.Worker](i), nil
}
