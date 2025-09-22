//go:build wireinject
// +build wireinject

package mealplangrocerylistinitializer

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/config"
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning/grocerylistpreparation"
	msgconfig "github.com/dinnerdonebetter/backend/internal/platform/messagequeue/config"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	loggingcfg "github.com/dinnerdonebetter/backend/internal/platform/observability/logging/config"
	metricscfg "github.com/dinnerdonebetter/backend/internal/platform/observability/metrics/config"
	tracingcfg "github.com/dinnerdonebetter/backend/internal/platform/observability/tracing/config"
	mealplangrocerylistinitializer "github.com/dinnerdonebetter/backend/internal/services/mealplanning/workers/meal_plan_grocery_list_initializer"

	"github.com/google/wire"
)

// Build builds a server.
func Build(
	ctx context.Context,
	cfg *config.MealPlanGroceryListInitializerConfig,
) (*mealplangrocerylistinitializer.Worker, error) {
	wire.Build(
		mealplangrocerylistinitializer.ProvidersMealPlanGroceryListInitializer,
		tracingcfg.TracingConfigProviders,
		observability.O11yProviders,
		msgconfig.MessageQueueProviders,
		loggingcfg.LogConfigProviders,
		metricscfg.MetricsConfigProviders,
		grocerylistpreparation.ProvidersGroceryListPreparation,
		ConfigProviders,
	)

	return nil, nil
}
