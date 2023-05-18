//go:build wireinject
// +build wireinject

package server

import (
	"context"

	analyticscfg "github.com/dinnerdonebetter/backend/internal/analytics/config"
	"github.com/dinnerdonebetter/backend/internal/authentication"
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
	"github.com/dinnerdonebetter/backend/internal/random"
	"github.com/dinnerdonebetter/backend/internal/routing/chi"
	"github.com/dinnerdonebetter/backend/internal/server"
	adminservice "github.com/dinnerdonebetter/backend/internal/services/admin"
	apiclientsservice "github.com/dinnerdonebetter/backend/internal/services/apiclients"
	authservice "github.com/dinnerdonebetter/backend/internal/services/authentication"
	householdinvitationssservice "github.com/dinnerdonebetter/backend/internal/services/householdinvitations"
	householdsservice "github.com/dinnerdonebetter/backend/internal/services/households"
	mealplaneventsservice "github.com/dinnerdonebetter/backend/internal/services/mealplanevents"
	"github.com/dinnerdonebetter/backend/internal/services/mealplangrocerylistitems"
	mealplanoptionsservice "github.com/dinnerdonebetter/backend/internal/services/mealplanoptions"
	mealplanoptionvotesservice "github.com/dinnerdonebetter/backend/internal/services/mealplanoptionvotes"
	mealplansservice "github.com/dinnerdonebetter/backend/internal/services/mealplans"
	"github.com/dinnerdonebetter/backend/internal/services/mealplantasks"
	mealsservice "github.com/dinnerdonebetter/backend/internal/services/meals"
	recipepreptasksservice "github.com/dinnerdonebetter/backend/internal/services/recipepreptasks"
	recipesservice "github.com/dinnerdonebetter/backend/internal/services/recipes"
	recipestepcompletionconditionsservice "github.com/dinnerdonebetter/backend/internal/services/recipestepcompletionconditions"
	recipestepingredientsservice "github.com/dinnerdonebetter/backend/internal/services/recipestepingredients"
	recipestepinstrumentsservice "github.com/dinnerdonebetter/backend/internal/services/recipestepinstruments"
	recipestepproductsservice "github.com/dinnerdonebetter/backend/internal/services/recipestepproducts"
	recipestepsservice "github.com/dinnerdonebetter/backend/internal/services/recipesteps"
	recipestepvesselsservice "github.com/dinnerdonebetter/backend/internal/services/recipestepvessels"
	"github.com/dinnerdonebetter/backend/internal/services/servicesettingconfigurations"
	servicesettingsservice "github.com/dinnerdonebetter/backend/internal/services/servicesettings"
	usersservice "github.com/dinnerdonebetter/backend/internal/services/users"
	validingredientmeasurementunitsservice "github.com/dinnerdonebetter/backend/internal/services/validingredientmeasurementunits"
	validingredientpreparationsservice "github.com/dinnerdonebetter/backend/internal/services/validingredientpreparations"
	validingredientsservice "github.com/dinnerdonebetter/backend/internal/services/validingredients"
	validingredientstateingredientsservice "github.com/dinnerdonebetter/backend/internal/services/validingredientstateingredients"
	validingredientstatesservice "github.com/dinnerdonebetter/backend/internal/services/validingredientstates"
	validinstrumentsservice "github.com/dinnerdonebetter/backend/internal/services/validinstruments"
	validmeasurementconversionsservice "github.com/dinnerdonebetter/backend/internal/services/validmeasurementconversions"
	"github.com/dinnerdonebetter/backend/internal/services/validmeasurementunits"
	validpreparationinstrumentsservice "github.com/dinnerdonebetter/backend/internal/services/validpreparationinstruments"
	validpreparationsservice "github.com/dinnerdonebetter/backend/internal/services/validpreparations"
	vendorproxyservice "github.com/dinnerdonebetter/backend/internal/services/vendorproxy"
	webhooksservice "github.com/dinnerdonebetter/backend/internal/services/webhooks"
	"github.com/dinnerdonebetter/backend/internal/uploads/images"

	"github.com/google/wire"
)

// Build builds a server.
func Build(
	ctx context.Context,
	cfg *config.InstanceConfig,
) (*server.HTTPServer, error) {
	wire.Build(
		config.ServiceConfigProviders,
		database.DBProviders,
		dbconfig.Providers,
		encoding.EncDecProviders,
		msgconfig.MessageQueueProviders,
		server.Providers,
		images.Providers,
		chi.Providers,
		random.Providers,
		authentication.Providers,
		authservice.Providers,
		usersservice.Providers,
		householdsservice.Providers,
		householdinvitationssservice.Providers,
		apiclientsservice.Providers,
		webhooksservice.Providers,
		adminservice.Providers,
		validinstrumentsservice.Providers,
		validingredientsservice.Providers,
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
		validmeasurementunits.Providers,
		validpreparationinstrumentsservice.Providers,
		validingredientstateingredientsservice.Providers,
		validingredientmeasurementunitsservice.Providers,
		mealplantasks.Providers,
		graphing.Providers,
		recipepreptasksservice.Providers,
		mealplangrocerylistitems.Providers,
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
		servicesettingconfigurations.Providers,
	)

	return nil, nil
}
