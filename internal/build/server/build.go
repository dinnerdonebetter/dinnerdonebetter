//go:build wireinject
// +build wireinject

package server

import (
	"context"

	analyticscfg "github.com/prixfixeco/backend/internal/analytics/config"
	"github.com/prixfixeco/backend/internal/authentication"
	"github.com/prixfixeco/backend/internal/config"
	"github.com/prixfixeco/backend/internal/database"
	dbconfig "github.com/prixfixeco/backend/internal/database/config"
	"github.com/prixfixeco/backend/internal/database/postgres"
	emailcfg "github.com/prixfixeco/backend/internal/email/config"
	"github.com/prixfixeco/backend/internal/encoding"
	featureflagscfg "github.com/prixfixeco/backend/internal/featureflags/config"
	graphing "github.com/prixfixeco/backend/internal/features/recipeanalysis"
	msgconfig "github.com/prixfixeco/backend/internal/messagequeue/config"
	"github.com/prixfixeco/backend/internal/observability"
	logcfg "github.com/prixfixeco/backend/internal/observability/logging/config"
	"github.com/prixfixeco/backend/internal/observability/tracing"
	tracingcfg "github.com/prixfixeco/backend/internal/observability/tracing/config"
	"github.com/prixfixeco/backend/internal/random"
	"github.com/prixfixeco/backend/internal/routing/chi"
	"github.com/prixfixeco/backend/internal/server"
	adminservice "github.com/prixfixeco/backend/internal/services/admin"
	apiclientsservice "github.com/prixfixeco/backend/internal/services/apiclients"
	authservice "github.com/prixfixeco/backend/internal/services/authentication"
	householdinvitationssservice "github.com/prixfixeco/backend/internal/services/householdinvitations"
	householdsservice "github.com/prixfixeco/backend/internal/services/households"
	mealplaneventsservice "github.com/prixfixeco/backend/internal/services/mealplanevents"
	"github.com/prixfixeco/backend/internal/services/mealplangrocerylistitems"
	mealplanoptionsservice "github.com/prixfixeco/backend/internal/services/mealplanoptions"
	mealplanoptionvotesservice "github.com/prixfixeco/backend/internal/services/mealplanoptionvotes"
	mealplansservice "github.com/prixfixeco/backend/internal/services/mealplans"
	"github.com/prixfixeco/backend/internal/services/mealplantasks"
	mealsservice "github.com/prixfixeco/backend/internal/services/meals"
	recipepreptasksservice "github.com/prixfixeco/backend/internal/services/recipepreptasks"
	recipesservice "github.com/prixfixeco/backend/internal/services/recipes"
	recipestepcompletionconditionsservice "github.com/prixfixeco/backend/internal/services/recipestepcompletionconditions"
	recipestepingredientsservice "github.com/prixfixeco/backend/internal/services/recipestepingredients"
	recipestepinstrumentsservice "github.com/prixfixeco/backend/internal/services/recipestepinstruments"
	recipestepproductsservice "github.com/prixfixeco/backend/internal/services/recipestepproducts"
	recipestepsservice "github.com/prixfixeco/backend/internal/services/recipesteps"
	recipestepvesselsservice "github.com/prixfixeco/backend/internal/services/recipestepvessels"
	"github.com/prixfixeco/backend/internal/services/servicesettingconfigurations"
	servicesettingsservice "github.com/prixfixeco/backend/internal/services/servicesettings"
	usersservice "github.com/prixfixeco/backend/internal/services/users"
	validingredientmeasurementunitsservice "github.com/prixfixeco/backend/internal/services/validingredientmeasurementunits"
	validingredientpreparationsservice "github.com/prixfixeco/backend/internal/services/validingredientpreparations"
	validingredientsservice "github.com/prixfixeco/backend/internal/services/validingredients"
	validingredientstateingredientsservice "github.com/prixfixeco/backend/internal/services/validingredientstateingredients"
	validingredientstatesservice "github.com/prixfixeco/backend/internal/services/validingredientstates"
	validinstrumentsservice "github.com/prixfixeco/backend/internal/services/validinstruments"
	validmeasurementconversionsservice "github.com/prixfixeco/backend/internal/services/validmeasurementconversions"
	"github.com/prixfixeco/backend/internal/services/validmeasurementunits"
	validpreparationinstrumentsservice "github.com/prixfixeco/backend/internal/services/validpreparationinstruments"
	validpreparationsservice "github.com/prixfixeco/backend/internal/services/validpreparations"
	vendorproxyservice "github.com/prixfixeco/backend/internal/services/vendorproxy"
	webhooksservice "github.com/prixfixeco/backend/internal/services/webhooks"
	"github.com/prixfixeco/backend/internal/uploads/images"

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
