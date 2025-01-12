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
	recipepreptasksservice "github.com/dinnerdonebetter/backend/internal/services/eating/recipepreptasks"
	reciperatingsservice "github.com/dinnerdonebetter/backend/internal/services/eating/reciperatings"
	recipesservice "github.com/dinnerdonebetter/backend/internal/services/eating/recipes"
	recipestepcompletionconditionsservice "github.com/dinnerdonebetter/backend/internal/services/eating/recipestepcompletionconditions"
	recipestepingredientsservice "github.com/dinnerdonebetter/backend/internal/services/eating/recipestepingredients"
	recipestepinstrumentsservice "github.com/dinnerdonebetter/backend/internal/services/eating/recipestepinstruments"
	recipestepproductsservice "github.com/dinnerdonebetter/backend/internal/services/eating/recipestepproducts"
	recipestepsservice "github.com/dinnerdonebetter/backend/internal/services/eating/recipesteps"
	recipestepvesselsservice "github.com/dinnerdonebetter/backend/internal/services/eating/recipestepvessels"
	useringredientpreferencesservice "github.com/dinnerdonebetter/backend/internal/services/eating/useringredientpreferences"
	validingredientgroupsservice "github.com/dinnerdonebetter/backend/internal/services/eating/validingredientgroups"
	validingredientmeasurementunitsservice "github.com/dinnerdonebetter/backend/internal/services/eating/validingredientmeasurementunits"
	validingredientpreparationsservice "github.com/dinnerdonebetter/backend/internal/services/eating/validingredientpreparations"
	validingredientsservice "github.com/dinnerdonebetter/backend/internal/services/eating/validingredients"
	validingredientstateingredientsservice "github.com/dinnerdonebetter/backend/internal/services/eating/validingredientstateingredients"
	validingredientstatesservice "github.com/dinnerdonebetter/backend/internal/services/eating/validingredientstates"
	validinstrumentsservice "github.com/dinnerdonebetter/backend/internal/services/eating/validinstruments"
	validmeasurementconversionsservice "github.com/dinnerdonebetter/backend/internal/services/eating/validmeasurementunitconversions"
	validmeasurementunitsservice "github.com/dinnerdonebetter/backend/internal/services/eating/validmeasurementunits"
	validpreparationinstrumentsservice "github.com/dinnerdonebetter/backend/internal/services/eating/validpreparationinstruments"
	validpreparationsservice "github.com/dinnerdonebetter/backend/internal/services/eating/validpreparations"
	validpreparationvesselsservice "github.com/dinnerdonebetter/backend/internal/services/eating/validpreparationvessels"
	validvesselsservice "github.com/dinnerdonebetter/backend/internal/services/eating/validvessels"
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
		validinstrumentsservice.Providers,
		validingredientsservice.Providers,
		validingredientgroupsservice.Providers,
		validpreparationsservice.Providers,
		validingredientpreparationsservice.Providers,
		mealsservice.Providers,
		recipesservice.Providers,
		recipestepsservice.Providers,
		recipestepproductsservice.Providers,
		recipestepinstrumentsservice.Providers,
		recipestepvesselsservice.Providers,
		recipestepingredientsservice.Providers,
		mealplansservice.Providers,
		mealplaneventsservice.Providers,
		mealplanoptionsservice.Providers,
		mealplanoptionvotesservice.Providers,
		validmeasurementunitsservice.Providers,
		validpreparationinstrumentsservice.Providers,
		validingredientstateingredientsservice.Providers,
		validingredientmeasurementunitsservice.Providers,
		mealplantasksservice.Providers,
		recipepreptasksservice.Providers,
		mealplangrocerylistitemsservice.Providers,
		validmeasurementconversionsservice.Providers,
		validingredientstatesservice.Providers,
		recipestepcompletionconditionsservice.Providers,
		servicesettingsservice.Providers,
		servicesettingconfigurationsservice.Providers,
		useringredientpreferencesservice.Providers,
		householdinstrumentownershipsservice.Providers,
		reciperatingsservice.Providers,
		oauth2clientsservice.Providers,
		validvesselsservice.Providers,
		validpreparationvesselsservice.Providers,
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
