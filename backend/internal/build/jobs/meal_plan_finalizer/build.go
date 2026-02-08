//go:build wireinject
// +build wireinject

package mealplanfinalizer

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/config"
	databasecfg "github.com/dinnerdonebetter/backend/internal/platform/database/config"
	"github.com/dinnerdonebetter/backend/internal/platform/database/postgres"
	msgconfig "github.com/dinnerdonebetter/backend/internal/platform/messagequeue/config"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	loggingcfg "github.com/dinnerdonebetter/backend/internal/platform/observability/logging/config"
	metricscfg "github.com/dinnerdonebetter/backend/internal/platform/observability/metrics/config"
	tracingcfg "github.com/dinnerdonebetter/backend/internal/platform/observability/tracing/config"
	"github.com/dinnerdonebetter/backend/internal/repositories/postgres/auditlogentries"
	"github.com/dinnerdonebetter/backend/internal/repositories/postgres/identity"
	"github.com/dinnerdonebetter/backend/internal/repositories/postgres/mealplanning"
	mealplanfinalizer "github.com/dinnerdonebetter/backend/internal/services/mealplanning/workers/meal_plan_finalizer"

	"github.com/google/wire"
)

// Build builds a server.
func Build(
	ctx context.Context,
	cfg *config.MealPlanFinalizerConfig,
) (*mealplanfinalizer.Worker, error) {
	wire.Build(
		mealplanfinalizer.ProvidersMealPlanFinalizer,
		tracingcfg.TracingConfigProviders,
		observability.O11yProviders,
		databasecfg.ClientConfigProviders,
		postgres.PGProviders,
		msgconfig.MessageQueueProviders,
		loggingcfg.LogConfigProviders,
		metricscfg.MetricsConfigProviders,
		auditlogentries.AuditRepoProviders,
		identity.IDRepoProviders,
		mealplanning.MPRepoProviders,
		ConfigProviders,
	)

	return nil, nil
}
