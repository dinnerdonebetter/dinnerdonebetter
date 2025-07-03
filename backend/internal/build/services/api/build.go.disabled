//go:build wireinject
// +build wireinject

package api

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/config"
	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/database/postgres"
	analyticscfg "github.com/dinnerdonebetter/backend/internal/platform/analytics/config"
	"github.com/dinnerdonebetter/backend/internal/platform/authentication"
	"github.com/dinnerdonebetter/backend/internal/platform/encoding"
	featureflagscfg "github.com/dinnerdonebetter/backend/internal/platform/featureflags/config"
	msgconfig "github.com/dinnerdonebetter/backend/internal/platform/messagequeue/config"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	loggingcfg "github.com/dinnerdonebetter/backend/internal/platform/observability/logging/config"
	metricscfg "github.com/dinnerdonebetter/backend/internal/platform/observability/metrics/config"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	tracingcfg "github.com/dinnerdonebetter/backend/internal/platform/observability/tracing/config"
	"github.com/dinnerdonebetter/backend/internal/platform/random"
	routingcfg "github.com/dinnerdonebetter/backend/internal/platform/routing/config"
	"github.com/dinnerdonebetter/backend/internal/platform/server/http"
	"github.com/dinnerdonebetter/backend/internal/platform/uploads/images"
	accountinvitationssservice "github.com/dinnerdonebetter/backend/internal/services/core/handlers/accountinvitations"
	accountsservice "github.com/dinnerdonebetter/backend/internal/services/core/handlers/accounts"
	adminservice "github.com/dinnerdonebetter/backend/internal/services/core/handlers/admin"
	auditlogentriesservice "github.com/dinnerdonebetter/backend/internal/services/core/handlers/auditlogentries"
	authservice "github.com/dinnerdonebetter/backend/internal/services/core/handlers/authentication"
	dataprivacyservice "github.com/dinnerdonebetter/backend/internal/services/core/handlers/dataprivacy"
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
		loggingcfg.ProvidersLogConfig,
		graphing.ProvidersRecipeAnalysis,
		authservice.Providers,
		usersservice.Providers,
		accountsservice.Providers,
		accountinvitationssservice.Providers,
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
