//go:build wireinject
// +build wireinject

package mealplanfinalizer

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/config"
	"github.com/dinnerdonebetter/backend/internal/database/postgres"
	msgconfig "github.com/dinnerdonebetter/backend/internal/lib/messagequeue/config"
	"github.com/dinnerdonebetter/backend/internal/lib/observability"
	loggingcfg "github.com/dinnerdonebetter/backend/internal/lib/observability/logging/config"
	metricscfg "github.com/dinnerdonebetter/backend/internal/lib/observability/metrics/config"
	tracingcfg "github.com/dinnerdonebetter/backend/internal/lib/observability/tracing/config"
	mealplanfinalizer "github.com/dinnerdonebetter/backend/internal/services/eating/workers/meal_plan_finalizer"

	"github.com/google/wire"
)

// Build builds a server.
func Build(
	ctx context.Context,
	cfg *config.MealPlanFinalizerConfig,
) (*mealplanfinalizer.Worker, error) {
	wire.Build(
		mealplanfinalizer.ProvidersMealPlanFinalizer,
		tracingcfg.ProvidersTracingConfig,
		observability.ProvidersObservability,
		postgres.ProvidersPostgres,
		msgconfig.MessageQueueProviders,
		loggingcfg.ProvidersLogConfig,
		metricscfg.ProvidersMetrics,
		ConfigProviders,
	)

	return nil, nil
}
