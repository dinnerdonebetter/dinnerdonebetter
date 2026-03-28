package mealplanfinalizer

import (
	"context"

	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/authorization"
	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/config"
	domaincustomroles "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/customroles"
	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/repositories/postgres/auditlogentries"
	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/repositories/postgres/customroles"
	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/repositories/postgres/identity"
	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/repositories/postgres/mealplanning"
	mealplanfinalizer "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/services/mealplanning/workers/meal_plan_finalizer"

	databasecfg "github.com/verygoodsoftwarenotvirus/platform/v4/database/config"
	"github.com/verygoodsoftwarenotvirus/platform/v4/database/postgres"
	msgconfig "github.com/verygoodsoftwarenotvirus/platform/v4/messagequeue/config"
	"github.com/verygoodsoftwarenotvirus/platform/v4/observability"
	loggingcfg "github.com/verygoodsoftwarenotvirus/platform/v4/observability/logging/config"
	metricscfg "github.com/verygoodsoftwarenotvirus/platform/v4/observability/metrics/config"
	tracingcfg "github.com/verygoodsoftwarenotvirus/platform/v4/observability/tracing/config"

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
	customroles.RegisterCustomRolesRepository(i)
	do.Provide(i, func(i do.Injector) (*authorization.RolePermissionCache, error) {
		cache := authorization.NewRolePermissionCache()
		repo := do.MustInvoke[domaincustomroles.Repository](i)
		if err := cache.Refresh(ctx, repo.GetAllCustomRolePermissions); err != nil {
			return nil, err
		}
		return cache, nil
	})
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
