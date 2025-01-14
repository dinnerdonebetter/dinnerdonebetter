//go:build wireinject
// +build wireinject

package api

import (
	"context"

	analyticscfg "github.com/dinnerdonebetter/backend/internal/analytics/config"
	"github.com/dinnerdonebetter/backend/internal/authentication"
	"github.com/dinnerdonebetter/backend/internal/config"
	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/database/postgres"
	"github.com/dinnerdonebetter/backend/internal/encoding"
	featureflagscfg "github.com/dinnerdonebetter/backend/internal/featureflags/config"
	graphing "github.com/dinnerdonebetter/backend/internal/features/recipeanalysis"
	msgconfig "github.com/dinnerdonebetter/backend/internal/messagequeue/config"
	"github.com/dinnerdonebetter/backend/internal/observability"
	loggingcfg "github.com/dinnerdonebetter/backend/internal/observability/logging/config"
	metricscfg "github.com/dinnerdonebetter/backend/internal/observability/metrics/config"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	tracingcfg "github.com/dinnerdonebetter/backend/internal/observability/tracing/config"
	"github.com/dinnerdonebetter/backend/internal/pkg/random"
	"github.com/dinnerdonebetter/backend/internal/routing/config"
	"github.com/dinnerdonebetter/backend/internal/server/http"
	adminservice "github.com/dinnerdonebetter/backend/internal/services/core/admin"
	auditlogentriesservice "github.com/dinnerdonebetter/backend/internal/services/core/auditlogentries"
	authservice "github.com/dinnerdonebetter/backend/internal/services/core/authentication"
	dataprivacyservice "github.com/dinnerdonebetter/backend/internal/services/core/dataprivacy"
	householdinvitationssservice "github.com/dinnerdonebetter/backend/internal/services/core/householdinvitations"
	householdsservice "github.com/dinnerdonebetter/backend/internal/services/core/households"
	oauth2clientsservice "github.com/dinnerdonebetter/backend/internal/services/core/oauth2clients"
	servicesettingconfigurationsservice "github.com/dinnerdonebetter/backend/internal/services/core/servicesettingconfigurations"
	servicesettingsservice "github.com/dinnerdonebetter/backend/internal/services/core/servicesettings"
	usernotificationsservice "github.com/dinnerdonebetter/backend/internal/services/core/usernotifications"
	usersservice "github.com/dinnerdonebetter/backend/internal/services/core/users"
	webhooksservice "github.com/dinnerdonebetter/backend/internal/services/core/webhooks"
	workersservice "github.com/dinnerdonebetter/backend/internal/services/core/workers"
	householdinstrumentownershipsservice "github.com/dinnerdonebetter/backend/internal/services/eating/householdinstrumentownerships"
	mealplaneventsservice "github.com/dinnerdonebetter/backend/internal/services/eating/mealplanevents"
	mealplangrocerylistitemsservice "github.com/dinnerdonebetter/backend/internal/services/eating/mealplangrocerylistitems"
	mealplanoptionsservice "github.com/dinnerdonebetter/backend/internal/services/eating/mealplanoptions"
	mealplanoptionvotesservice "github.com/dinnerdonebetter/backend/internal/services/eating/mealplanoptionvotes"
	mealplansservice "github.com/dinnerdonebetter/backend/internal/services/eating/mealplans"
	mealplantasksservice "github.com/dinnerdonebetter/backend/internal/services/eating/mealplantasks"
	mealsservice "github.com/dinnerdonebetter/backend/internal/services/eating/meals"
	recipemanagementservice "github.com/dinnerdonebetter/backend/internal/services/eating/recipe_management"
	useringredientpreferencesservice "github.com/dinnerdonebetter/backend/internal/services/eating/useringredientpreferences"
	validenumerationsservice "github.com/dinnerdonebetter/backend/internal/services/eating/valid_enumerations"
	"github.com/dinnerdonebetter/backend/internal/uploads/images"

	"github.com/google/wire"
)

// Build builds a server.
func Build(
	ctx context.Context,
	cfg *config.APIServiceConfig,
) (http.Server, error) {
	wire.Build(
		authentication.AuthProviders,
		config.ServiceConfigProviders,
		database.DBProviders,
		encoding.EncDecProviders,
		msgconfig.MessageQueueProviders,
		http.ProvidersHTTP,
		images.ProvidersImages,
		random.ProvidersRandom,
		featureflagscfg.ProvidersFeatureFlags,
		tracing.ProvidersTracing,
		tracingcfg.ProvidersTracing,
		observability.ProvidersObservability,
		postgres.ProvidersPostgres,
		loggingcfg.ProvidersLogConfig,
		graphing.Providers,
		authservice.Providers,
		usersservice.Providers,
		householdsservice.Providers,
		householdinvitationssservice.Providers,
		webhooksservice.Providers,
		adminservice.Providers,
		validenumerationsservice.Providers,
		mealsservice.Providers,
		recipemanagementservice.Providers,
		mealplansservice.Providers,
		mealplaneventsservice.Providers,
		mealplanoptionsservice.Providers,
		mealplanoptionvotesservice.Providers,
		mealplantasksservice.Providers,
		mealplangrocerylistitemsservice.Providers,
		servicesettingsservice.Providers,
		servicesettingconfigurationsservice.Providers,
		useringredientpreferencesservice.Providers,
		householdinstrumentownershipsservice.Providers,
		oauth2clientsservice.Providers,
		analyticscfg.ProvidersAnalytics,
		workersservice.Providers,
		usernotificationsservice.Providers,
		auditlogentriesservice.Providers,
		dataprivacyservice.Providers,
		metricscfg.ProvidersMetrics,
		routingcfg.RoutingConfigProviders,
	)

	return nil, nil
}
