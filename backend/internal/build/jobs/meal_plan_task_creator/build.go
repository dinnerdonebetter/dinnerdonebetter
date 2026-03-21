package mealplantaskcreator

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/config"
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning/recipeanalysis"
	"github.com/dinnerdonebetter/backend/internal/repositories/postgres/auditlogentries"
	"github.com/dinnerdonebetter/backend/internal/repositories/postgres/identity"
	"github.com/dinnerdonebetter/backend/internal/repositories/postgres/mealplanning"
	mealplantaskcreator "github.com/dinnerdonebetter/backend/internal/services/mealplanning/workers/meal_plan_task_creator"

	"github.com/samber/do/v2"
	databasecfg "github.com/verygoodsoftwarenotvirus/platform/database/config"
	"github.com/verygoodsoftwarenotvirus/platform/database/postgres"
	msgconfig "github.com/verygoodsoftwarenotvirus/platform/messagequeue/config"
	"github.com/verygoodsoftwarenotvirus/platform/observability"
	loggingcfg "github.com/verygoodsoftwarenotvirus/platform/observability/logging/config"
	metricscfg "github.com/verygoodsoftwarenotvirus/platform/observability/metrics/config"
	tracingcfg "github.com/verygoodsoftwarenotvirus/platform/observability/tracing/config"
)

// BuildInjector creates and configures the dependency injection container.
func BuildInjector(
	ctx context.Context,
	cfg *config.MealPlanTaskCreatorConfig,
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
	recipeanalysis.RegisterRecipeAnalyzer(i)
	auditlogentries.RegisterAuditLogRepository(i)
	identity.RegisterIdentityRepository(i)
	mealplanning.RegisterMealPlanningRepository(i)
	mealplantaskcreator.RegisterMealPlanTaskCreator(i)

	return i
}

// Build builds a server.
func Build(
	ctx context.Context,
	cfg *config.MealPlanTaskCreatorConfig,
) (*mealplantaskcreator.Worker, error) {
	i := BuildInjector(ctx, cfg)
	return do.MustInvoke[*mealplantaskcreator.Worker](i), nil
}
