//go:build wireinject
// +build wireinject

package mealplanfinalizer

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/config"
	"github.com/dinnerdonebetter/backend/internal/platform/database/postgres"
	msgconfig "github.com/dinnerdonebetter/backend/internal/platform/messagequeue/config"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	loggingcfg "github.com/dinnerdonebetter/backend/internal/platform/observability/logging/config"
	metricscfg "github.com/dinnerdonebetter/backend/internal/platform/observability/metrics/config"
	tracingcfg "github.com/dinnerdonebetter/backend/internal/platform/observability/tracing/config"
	"github.com/dinnerdonebetter/backend/internal/repositories/postgres/auditlogentries"
	"github.com/dinnerdonebetter/backend/internal/repositories/postgres/identity"
	"github.com/dinnerdonebetter/backend/internal/repositories/postgres/mealplanning"
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
		observability.Providers,
		postgres.Providers,
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
