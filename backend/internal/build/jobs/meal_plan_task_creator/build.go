//go:build wireinject
// +build wireinject

package mealplantaskcreator

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/config"
	"github.com/dinnerdonebetter/backend/internal/database/postgres"
	msgconfig "github.com/dinnerdonebetter/backend/internal/lib/messagequeue/config"
	"github.com/dinnerdonebetter/backend/internal/lib/observability"
	loggingcfg "github.com/dinnerdonebetter/backend/internal/lib/observability/logging/config"
	metricscfg "github.com/dinnerdonebetter/backend/internal/lib/observability/metrics/config"
	tracingcfg "github.com/dinnerdonebetter/backend/internal/lib/observability/tracing/config"
	"github.com/dinnerdonebetter/backend/internal/services/eating/businesslogic/recipeanalysis"
	mealplantaskcreator "github.com/dinnerdonebetter/backend/internal/services/eating/workers/meal_plan_task_creator"

	"github.com/google/wire"
)

// Build builds a server.
func Build(
	ctx context.Context,
	cfg *config.MealPlanTaskCreatorConfig,
) (*mealplantaskcreator.Worker, error) {
	wire.Build(
		postgres.ProvidersPostgres,
		recipeanalysis.ProvidersRecipeAnalysis,
		mealplantaskcreator.ProvidersMealPlanTaskCreator,
		tracingcfg.ProvidersTracingConfig,
		observability.ProvidersObservability,
		msgconfig.MessageQueueProviders,
		loggingcfg.ProvidersLogConfig,
		metricscfg.ProvidersMetrics,
		ConfigProviders,
	)

	return nil, nil
}
