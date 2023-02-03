//go:build wireinject
// +build wireinject

package server

import (
	"context"

	"github.com/google/wire"

	"github.com/prixfixeco/backend/internal/authentication"
	"github.com/prixfixeco/backend/internal/config"
	"github.com/prixfixeco/backend/internal/database"
	dbconfig "github.com/prixfixeco/backend/internal/database/config"
	"github.com/prixfixeco/backend/internal/email"
	"github.com/prixfixeco/backend/internal/encoding"
	graphing "github.com/prixfixeco/backend/internal/features/recipeanalysis"
	msgconfig "github.com/prixfixeco/backend/internal/messagequeue/config"
	"github.com/prixfixeco/backend/internal/observability/logging"
	"github.com/prixfixeco/backend/internal/observability/metrics"
	"github.com/prixfixeco/backend/internal/observability/tracing"
	"github.com/prixfixeco/backend/internal/random"
	"github.com/prixfixeco/backend/internal/routing/chi"
	"github.com/prixfixeco/backend/internal/server"
	adminservice "github.com/prixfixeco/backend/internal/services/admin"
	apiclientsservice "github.com/prixfixeco/backend/internal/services/apiclients"
	authservice "github.com/prixfixeco/backend/internal/services/authentication"
	householdinvitationssservice "github.com/prixfixeco/backend/internal/services/householdinvitations"
	householdsservice "github.com/prixfixeco/backend/internal/services/households"
	mealplaneventsservice "github.com/prixfixeco/backend/internal/services/mealplanevents"
	mealplangrocerylistitems "github.com/prixfixeco/backend/internal/services/mealplangrocerylistitems"
	mealplanoptionsservice "github.com/prixfixeco/backend/internal/services/mealplanoptions"
	mealplanoptionvotesservice "github.com/prixfixeco/backend/internal/services/mealplanoptionvotes"
	mealplansservice "github.com/prixfixeco/backend/internal/services/mealplans"
	mealplantasks "github.com/prixfixeco/backend/internal/services/mealplantasks"
	mealsservice "github.com/prixfixeco/backend/internal/services/meals"
	recipepreptasksservice "github.com/prixfixeco/backend/internal/services/recipepreptasks"
	recipesservice "github.com/prixfixeco/backend/internal/services/recipes"
	recipestepingredientsservice "github.com/prixfixeco/backend/internal/services/recipestepingredients"
	recipestepinstrumentsservice "github.com/prixfixeco/backend/internal/services/recipestepinstruments"
	recipestepproductsservice "github.com/prixfixeco/backend/internal/services/recipestepproducts"
	recipestepsservice "github.com/prixfixeco/backend/internal/services/recipesteps"
	usersservice "github.com/prixfixeco/backend/internal/services/users"
	validingredientmeasurementunitsservice "github.com/prixfixeco/backend/internal/services/validingredientmeasurementunits"
	validingredientpreparationsservice "github.com/prixfixeco/backend/internal/services/validingredientpreparations"
	validingredientsservice "github.com/prixfixeco/backend/internal/services/validingredients"
	validingredientstateingredientsservice "github.com/prixfixeco/backend/internal/services/validingredientstateingredients"
	validingredientstatesservice "github.com/prixfixeco/backend/internal/services/validingredientstates"
	validinstrumentsservice "github.com/prixfixeco/backend/internal/services/validinstruments"
	validmeasurementconversionsservice "github.com/prixfixeco/backend/internal/services/validmeasurementconversions"
	validmeasurementunits "github.com/prixfixeco/backend/internal/services/validmeasurementunits"
	validpreparationinstrumentsservice "github.com/prixfixeco/backend/internal/services/validpreparationinstruments"
	validpreparationsservice "github.com/prixfixeco/backend/internal/services/validpreparations"
	webhooksservice "github.com/prixfixeco/backend/internal/services/webhooks"
	"github.com/prixfixeco/backend/internal/uploads/images"
)

// Build builds a server.
func Build(
	ctx context.Context,
	logger logging.Logger,
	cfg *config.InstanceConfig,
	tracerProvider tracing.TracerProvider,
	unitCounterProvider metrics.UnitCounterProvider,
	metricsHandler metrics.Handler,
	dataManager database.DataManager,
	emailer email.Emailer,
) (*server.HTTPServer, error) {
	wire.Build(
		config.Providers,
		database.Providers,
		dbconfig.Providers,
		encoding.Providers,
		msgconfig.Providers,
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
	)

	return nil, nil
}
