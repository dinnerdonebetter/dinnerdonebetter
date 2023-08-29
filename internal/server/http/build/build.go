//go:build wireinject
// +build wireinject

package build

import (
	"context"

	analyticscfg "github.com/dinnerdonebetter/backend/internal/analytics/config"
	authcfg "github.com/dinnerdonebetter/backend/internal/authentication/config"
	"github.com/dinnerdonebetter/backend/internal/config"
	"github.com/dinnerdonebetter/backend/internal/database"
	dbconfig "github.com/dinnerdonebetter/backend/internal/database/config"
	"github.com/dinnerdonebetter/backend/internal/database/postgres"
	emailcfg "github.com/dinnerdonebetter/backend/internal/email/config"
	"github.com/dinnerdonebetter/backend/internal/encoding"
	featureflagscfg "github.com/dinnerdonebetter/backend/internal/featureflags/config"
	graphing "github.com/dinnerdonebetter/backend/internal/features/recipeanalysis"
	msgconfig "github.com/dinnerdonebetter/backend/internal/messagequeue/config"
	"github.com/dinnerdonebetter/backend/internal/observability"
	logcfg "github.com/dinnerdonebetter/backend/internal/observability/logging/config"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	tracingcfg "github.com/dinnerdonebetter/backend/internal/observability/tracing/config"
	"github.com/dinnerdonebetter/backend/internal/pkg/random"
	"github.com/dinnerdonebetter/backend/internal/routing/chi"
	"github.com/dinnerdonebetter/backend/internal/server/http"
	adminservice "github.com/dinnerdonebetter/backend/internal/services/admin"
	authservice "github.com/dinnerdonebetter/backend/internal/services/authentication"
	householdinstrumentownershipsservice "github.com/dinnerdonebetter/backend/internal/services/householdinstrumentownerships"
	householdinvitationssservice "github.com/dinnerdonebetter/backend/internal/services/householdinvitations"
	householdsservice "github.com/dinnerdonebetter/backend/internal/services/households"
	mealplaneventsservice "github.com/dinnerdonebetter/backend/internal/services/mealplanevents"
	mealplangrocerylistitemsservice "github.com/dinnerdonebetter/backend/internal/services/mealplangrocerylistitems"
	mealplanoptionsservice "github.com/dinnerdonebetter/backend/internal/services/mealplanoptions"
	mealplanoptionvotesservice "github.com/dinnerdonebetter/backend/internal/services/mealplanoptionvotes"
	mealplansservice "github.com/dinnerdonebetter/backend/internal/services/mealplans"
	mealplantasksservice "github.com/dinnerdonebetter/backend/internal/services/mealplantasks"
	mealsservice "github.com/dinnerdonebetter/backend/internal/services/meals"
	oauth2clientsservice "github.com/dinnerdonebetter/backend/internal/services/oauth2clients"
	recipepreptasksservice "github.com/dinnerdonebetter/backend/internal/services/recipepreptasks"
	reciperatingsservice "github.com/dinnerdonebetter/backend/internal/services/reciperatings"
	recipesservice "github.com/dinnerdonebetter/backend/internal/services/recipes"
	recipestepcompletionconditionsservice "github.com/dinnerdonebetter/backend/internal/services/recipestepcompletionconditions"
	recipestepingredientsservice "github.com/dinnerdonebetter/backend/internal/services/recipestepingredients"
	recipestepinstrumentsservice "github.com/dinnerdonebetter/backend/internal/services/recipestepinstruments"
	recipestepproductsservice "github.com/dinnerdonebetter/backend/internal/services/recipestepproducts"
	recipestepsservice "github.com/dinnerdonebetter/backend/internal/services/recipesteps"
	recipestepvesselsservice "github.com/dinnerdonebetter/backend/internal/services/recipestepvessels"
	servicesettingconfigurationsservice "github.com/dinnerdonebetter/backend/internal/services/servicesettingconfigurations"
	servicesettingsservice "github.com/dinnerdonebetter/backend/internal/services/servicesettings"
	useringredientpreferencesservice "github.com/dinnerdonebetter/backend/internal/services/useringredientpreferences"
	usersservice "github.com/dinnerdonebetter/backend/internal/services/users"
	validingredientgroupsservice "github.com/dinnerdonebetter/backend/internal/services/validingredientgroups"
	validingredientmeasurementunitsservice "github.com/dinnerdonebetter/backend/internal/services/validingredientmeasurementunits"
	validingredientpreparationsservice "github.com/dinnerdonebetter/backend/internal/services/validingredientpreparations"
	validingredientsservice "github.com/dinnerdonebetter/backend/internal/services/validingredients"
	validingredientstateingredientsservice "github.com/dinnerdonebetter/backend/internal/services/validingredientstateingredients"
	validingredientstatesservice "github.com/dinnerdonebetter/backend/internal/services/validingredientstates"
	validinstrumentsservice "github.com/dinnerdonebetter/backend/internal/services/validinstruments"
	validmeasurementconversionsservice "github.com/dinnerdonebetter/backend/internal/services/validmeasurementconversions"
	validmeasurementunitsservice "github.com/dinnerdonebetter/backend/internal/services/validmeasurementunits"
	validpreparationinstrumentsservice "github.com/dinnerdonebetter/backend/internal/services/validpreparationinstruments"
	validpreparationsservice "github.com/dinnerdonebetter/backend/internal/services/validpreparations"
	validpreparationvesselsservice "github.com/dinnerdonebetter/backend/internal/services/validpreparationvessels"
	validvesselsservice "github.com/dinnerdonebetter/backend/internal/services/validvessels"
	vendorproxyservice "github.com/dinnerdonebetter/backend/internal/services/vendorproxy"
	webhooksservice "github.com/dinnerdonebetter/backend/internal/services/webhooks"
	"github.com/dinnerdonebetter/backend/internal/uploads/images"

	"github.com/google/wire"
)

// Build builds a server.
func Build(
	ctx context.Context,
	cfg *config.InstanceConfig,
) (http.Server, error) {
	wire.Build(
		config.ServiceConfigProviders,
		database.DBProviders,
		dbconfig.Providers,
		encoding.EncDecProviders,
		msgconfig.MessageQueueProviders,
		http.Providers,
		images.Providers,
		chi.Providers,
		random.Providers,
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
		graphing.Providers,
		recipepreptasksservice.Providers,
		mealplangrocerylistitemsservice.Providers,
		validmeasurementconversionsservice.Providers,
		validingredientstatesservice.Providers,
		recipestepcompletionconditionsservice.Providers,
		featureflagscfg.Providers,
		vendorproxyservice.Providers,
		tracing.Providers,
		emailcfg.Providers,
		tracingcfg.Providers,
		observability.Providers,
		postgres.Providers,
		analyticscfg.Providers,
		logcfg.Providers,
		servicesettingsservice.Providers,
		servicesettingconfigurationsservice.Providers,
		useringredientpreferencesservice.Providers,
		householdinstrumentownershipsservice.Providers,
		reciperatingsservice.Providers,
		oauth2clientsservice.Providers,
		validvesselsservice.Providers,
		validpreparationvesselsservice.Providers,
		authcfg.Providers,
	)

	return nil, nil
}
