//go:build wireinject
// +build wireinject

package mealplangrocerylistinitializer

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/config"
	msgconfig "github.com/dinnerdonebetter/backend/internal/lib/messagequeue/config"
	"github.com/dinnerdonebetter/backend/internal/lib/observability"
	loggingcfg "github.com/dinnerdonebetter/backend/internal/lib/observability/logging/config"
	metricscfg "github.com/dinnerdonebetter/backend/internal/lib/observability/metrics/config"
	tracingcfg "github.com/dinnerdonebetter/backend/internal/lib/observability/tracing/config"
	mealplangrocerylistinitializer "github.com/dinnerdonebetter/backend/internal/services/eating/workers/meal_plan_grocery_list_initializer"

	"github.com/google/wire"
)

// Build builds a server.
func Build(
	ctx context.Context,
	cfg *config.MealPlanGroceryListInitializerConfig,
) (*mealplangrocerylistinitializer.Worker, error) {
	wire.Build(
		mealplangrocerylistinitializer.ProvidersMealPlanGroceryListInitializer,
		tracingcfg.ProvidersTracingConfig,
		observability.ProvidersObservability,
		msgconfig.MessageQueueProviders,
		loggingcfg.ProvidersLogConfig,
		metricscfg.ProvidersMetrics,
		// TODO: grocerylistpreparation.ProvidersGroceryListPreparation,
		ConfigProviders,
	)

	return nil, nil
}
