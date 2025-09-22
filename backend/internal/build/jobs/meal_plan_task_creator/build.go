//go:build wireinject
// +build wireinject

package mealplantaskcreator

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/config"
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning/recipeanalysis"
	"github.com/dinnerdonebetter/backend/internal/platform/database/postgres"
	msgconfig "github.com/dinnerdonebetter/backend/internal/platform/messagequeue/config"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	loggingcfg "github.com/dinnerdonebetter/backend/internal/platform/observability/logging/config"
	metricscfg "github.com/dinnerdonebetter/backend/internal/platform/observability/metrics/config"
	tracingcfg "github.com/dinnerdonebetter/backend/internal/platform/observability/tracing/config"
	"github.com/dinnerdonebetter/backend/internal/repositories/postgres/auditlogentries"
	"github.com/dinnerdonebetter/backend/internal/repositories/postgres/identity"
	"github.com/dinnerdonebetter/backend/internal/repositories/postgres/mealplanning"
	mealplantaskcreator "github.com/dinnerdonebetter/backend/internal/services/mealplanning/workers/meal_plan_task_creator"

	"github.com/google/wire"
)

// Build builds a server.
func Build(
	ctx context.Context,
	cfg *config.MealPlanTaskCreatorConfig,
) (*mealplantaskcreator.Worker, error) {
	wire.Build(
		postgres.Providers,
		recipeanalysis.ProvidersRecipeAnalysis,
		mealplantaskcreator.ProvidersMealPlanTaskCreator,
		tracingcfg.ProvidersTracingConfig,
		observability.Providers,
		msgconfig.MessageQueueProviders,
		loggingcfg.ProvidersLogConfig,
		metricscfg.Providers,
		auditlogentries.Providers,
		identity.Providers,
		mealplanning.Providers,
		ConfigProviders,
	)

	return nil, nil
}
