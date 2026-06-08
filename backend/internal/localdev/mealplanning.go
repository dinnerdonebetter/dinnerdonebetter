package localdev

import (
	"context"
	"fmt"

	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/repositories/postgres/auditlogentries"
	identityrepo "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/repositories/postgres/identity"
	mealplanningrepo "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/repositories/postgres/mealplanning"
	"github.com/dinnerdonebetter/dinnerdonebetter/backend/pkg/client"

	"github.com/primandproper/platform/database"
	databasecfg "github.com/primandproper/platform/database/config"
	"github.com/primandproper/platform/observability/logging"
	"github.com/primandproper/platform/observability/tracing"
)

// WithMealPlanningRepository provides a meal planning repository for custom operations.
// The provided function receives a fully configured mealplanning.Repository along with logger and tracer.
// This repository handles all meal planning entities including recipes, ingredients, preparations, vessels, instruments, etc.
func WithMealPlanningRepository(fn func(ctx context.Context, repo mealplanning.Repository, logger logging.Logger, tracerProvider tracing.TracerProvider) error) DatabaseInitFunc {
	return func(ctx context.Context, dbClient database.Client, dbCfg *databasecfg.Config, logger logging.Logger, tracerProvider tracing.TracerProvider) error {
		auditLogRepo := auditlogentries.ProvideAuditLogRepository(logger, tracerProvider, dbClient)
		identityRepo := identityrepo.ProvideIdentityRepository(logger, tracerProvider, auditLogRepo, dbClient)
		mealPlanningRepo := mealplanningrepo.ProvideMealPlanningRepository(logger, tracerProvider, auditLogRepo, identityRepo, dbClient)
		return fn(ctx, mealPlanningRepo, logger, tracerProvider)
	}
}

// BuildInsecureOAuthedMealPlanningGRPCClient performs the OAuth2 authorization-code flow and
// returns a MealPlanningClient. Use this in meal-planning integration tests and local dev tooling
// that needs to call meal-planning RPCs in addition to the base platform RPCs.
func BuildInsecureOAuthedMealPlanningGRPCClient(
	ctx context.Context,
	createdClientID,
	createdClientSecret,
	httpTestServerAddress,
	grpcServerAddress,
	token string,
) (client.MealPlanningClient, error) {
	opts, err := buildInsecureOAuthedDialOptions(ctx, createdClientID, createdClientSecret, httpTestServerAddress, token)
	if err != nil {
		return nil, err
	}

	c, err := client.BuildMealPlanningClient(grpcServerAddress, opts...)
	if err != nil {
		return nil, fmt.Errorf("building meal planning client: %w", err)
	}

	return c, nil
}
