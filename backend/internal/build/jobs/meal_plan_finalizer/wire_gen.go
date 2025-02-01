// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package mealplanfinalizer

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/config"
	"github.com/dinnerdonebetter/backend/internal/database/postgres"
	"github.com/dinnerdonebetter/backend/internal/lib/messagequeue/config"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/logging/config"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/metrics/config"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing/config"
	"github.com/dinnerdonebetter/backend/internal/services/eating/workers/meal_plan_finalizer"
)

// Injectors from build.go:

// Build builds a server.
func Build(ctx context.Context, cfg *config.MealPlanFinalizerConfig) (*mealplanfinalizer.Worker, error) {
	observabilityConfig := &cfg.Observability
	loggingcfgConfig := &observabilityConfig.Logging
	logger, err := loggingcfg.ProvideLogger(ctx, loggingcfgConfig)
	if err != nil {
		return nil, err
	}
	tracingcfgConfig := &observabilityConfig.Tracing
	tracerProvider, err := tracingcfg.ProvideTracerProvider(ctx, tracingcfgConfig, logger)
	if err != nil {
		return nil, err
	}
	databasecfgConfig := &cfg.Database
	dataManager, err := postgres.ProvideDatabaseClient(ctx, logger, tracerProvider, databasecfgConfig)
	if err != nil {
		return nil, err
	}
	msgconfigConfig := &cfg.Events
	publisherProvider, err := msgconfig.ProvidePublisherProvider(ctx, logger, tracerProvider, msgconfigConfig)
	if err != nil {
		return nil, err
	}
	metricscfgConfig := &observabilityConfig.Metrics
	provider, err := metricscfg.ProvideMetricsProvider(ctx, logger, metricscfgConfig)
	if err != nil {
		return nil, err
	}
	queuesConfig := cfg.Queues
	worker, err := mealplanfinalizer.NewMealPlanFinalizer(logger, tracerProvider, dataManager, publisherProvider, provider, queuesConfig)
	if err != nil {
		return nil, err
	}
	return worker, nil
}
