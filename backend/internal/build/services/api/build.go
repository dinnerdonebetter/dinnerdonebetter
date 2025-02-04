//go:build wireinject
// +build wireinject

package api

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/config"
	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/database/postgres"
	analyticscfg "github.com/dinnerdonebetter/backend/internal/lib/analytics/config"
	"github.com/dinnerdonebetter/backend/internal/lib/authentication"
	"github.com/dinnerdonebetter/backend/internal/lib/encoding"
	featureflagscfg "github.com/dinnerdonebetter/backend/internal/lib/featureflags/config"
	msgconfig "github.com/dinnerdonebetter/backend/internal/lib/messagequeue/config"
	"github.com/dinnerdonebetter/backend/internal/lib/observability"
	loggingcfg "github.com/dinnerdonebetter/backend/internal/lib/observability/logging/config"
	metricscfg "github.com/dinnerdonebetter/backend/internal/lib/observability/metrics/config"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
	tracingcfg "github.com/dinnerdonebetter/backend/internal/lib/observability/tracing/config"
	"github.com/dinnerdonebetter/backend/internal/lib/random"
	routingcfg "github.com/dinnerdonebetter/backend/internal/lib/routing/config"
	"github.com/dinnerdonebetter/backend/internal/lib/server/http"
	"github.com/dinnerdonebetter/backend/internal/lib/uploads/images"
	adminservice "github.com/dinnerdonebetter/backend/internal/services/core/handlers/admin"
	auditlogentriesservice "github.com/dinnerdonebetter/backend/internal/services/core/handlers/auditlogentries"
	authservice "github.com/dinnerdonebetter/backend/internal/services/core/handlers/authentication"
	dataprivacyservice "github.com/dinnerdonebetter/backend/internal/services/core/handlers/dataprivacy"
	householdinvitationssservice "github.com/dinnerdonebetter/backend/internal/services/core/handlers/householdinvitations"
	householdsservice "github.com/dinnerdonebetter/backend/internal/services/core/handlers/households"
	oauth2clientsservice "github.com/dinnerdonebetter/backend/internal/services/core/handlers/oauth2clients"
	servicesettingconfigurationsservice "github.com/dinnerdonebetter/backend/internal/services/core/handlers/servicesettingconfigurations"
	servicesettingsservice "github.com/dinnerdonebetter/backend/internal/services/core/handlers/servicesettings"
	usernotificationsservice "github.com/dinnerdonebetter/backend/internal/services/core/handlers/usernotifications"
	usersservice "github.com/dinnerdonebetter/backend/internal/services/core/handlers/users"
	webhooksservice "github.com/dinnerdonebetter/backend/internal/services/core/handlers/webhooks"
	"github.com/dinnerdonebetter/backend/internal/services/eating/businesslogic/grocerylistpreparation"
	graphing "github.com/dinnerdonebetter/backend/internal/services/eating/businesslogic/recipeanalysis"
	mealplanningservice "github.com/dinnerdonebetter/backend/internal/services/eating/handlers/meal_planning"
	recipemanagementservice "github.com/dinnerdonebetter/backend/internal/services/eating/handlers/recipe_management"
	validenumerationsservice "github.com/dinnerdonebetter/backend/internal/services/eating/handlers/valid_enumerations"
	workersservice "github.com/dinnerdonebetter/backend/internal/services/eating/handlers/workers"
	mealplanfinalizer "github.com/dinnerdonebetter/backend/internal/services/eating/workers/meal_plan_finalizer"
	mealplangrocerylistinitializer "github.com/dinnerdonebetter/backend/internal/services/eating/workers/meal_plan_grocery_list_initializer"
	mealplantaskcreator "github.com/dinnerdonebetter/backend/internal/services/eating/workers/meal_plan_task_creator"

	"github.com/google/wire"
)

// Build builds a server.
func Build(
	ctx context.Context,
	cfg *config.APIServiceConfig,
) (http.Server, error) {
	wire.Build(
		authentication.AuthProviders,
		ConfigProviders,
		database.DBProviders,
		encoding.EncDecProviders,
		msgconfig.MessageQueueProviders,
		images.ProvidersImages,
		random.ProvidersRandom,
		featureflagscfg.ProvidersFeatureFlags,
		tracing.ProvidersTracing,
		tracingcfg.ProvidersTracingConfig,
		observability.ProvidersObservability,
		postgres.ProvidersPostgres,
		loggingcfg.ProvidersLoggingConfig,
		graphing.ProvidersRecipeAnalysis,
		authservice.Providers,
		usersservice.Providers,
		householdsservice.Providers,
		householdinvitationssservice.Providers,
		webhooksservice.Providers,
		adminservice.Providers,
		validenumerationsservice.Providers,
		recipemanagementservice.Providers,
		servicesettingsservice.Providers,
		servicesettingconfigurationsservice.Providers,
		oauth2clientsservice.Providers,
		analyticscfg.ProvidersAnalytics,
		workersservice.Providers,
		usernotificationsservice.Providers,
		auditlogentriesservice.Providers,
		dataprivacyservice.Providers,
		metricscfg.ProvidersMetrics,
		mealplanningservice.Providers,
		http.ProvidersHTTP,
		routingcfg.RoutingConfigProviders,
		mealplantaskcreator.ProvidersMealPlanTaskCreator,
		mealplangrocerylistinitializer.ProvidersMealPlanGroceryListInitializer,
		mealplanfinalizer.ProvidersMealPlanFinalizer,
		grocerylistpreparation.ProvidersGroceryListPreparation,
		ProvideAPIRouter,
	)

	return nil, nil
}
