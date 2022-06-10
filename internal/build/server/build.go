//go:build wireinject
// +build wireinject

package server

import (
	"context"

	"github.com/google/wire"

	"github.com/prixfixeco/api_server/internal/authentication"
	"github.com/prixfixeco/api_server/internal/config"
	"github.com/prixfixeco/api_server/internal/database"
	dbconfig "github.com/prixfixeco/api_server/internal/database/config"
	"github.com/prixfixeco/api_server/internal/email"
	"github.com/prixfixeco/api_server/internal/encoding"
	msgconfig "github.com/prixfixeco/api_server/internal/messagequeue/config"
	"github.com/prixfixeco/api_server/internal/observability/logging"
	"github.com/prixfixeco/api_server/internal/observability/metrics"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/internal/routing/chi"
	"github.com/prixfixeco/api_server/internal/server"
	adminservice "github.com/prixfixeco/api_server/internal/services/admin"
	apiclientsservice "github.com/prixfixeco/api_server/internal/services/apiclients"
	authservice "github.com/prixfixeco/api_server/internal/services/authentication"
	householdinvitationssservice "github.com/prixfixeco/api_server/internal/services/householdinvitations"
	householdsservice "github.com/prixfixeco/api_server/internal/services/households"
	mealplanoptionsservice "github.com/prixfixeco/api_server/internal/services/mealplanoptions"
	mealplanoptionvotesservice "github.com/prixfixeco/api_server/internal/services/mealplanoptionvotes"
	mealplansservice "github.com/prixfixeco/api_server/internal/services/mealplans"
	mealsservice "github.com/prixfixeco/api_server/internal/services/meals"
	recipesservice "github.com/prixfixeco/api_server/internal/services/recipes"
	recipestepingredientsservice "github.com/prixfixeco/api_server/internal/services/recipestepingredients"
	recipestepinstrumentsservice "github.com/prixfixeco/api_server/internal/services/recipestepinstruments"
	recipestepproductsservice "github.com/prixfixeco/api_server/internal/services/recipestepproducts"
	recipestepsservice "github.com/prixfixeco/api_server/internal/services/recipesteps"
	usersservice "github.com/prixfixeco/api_server/internal/services/users"
	validingredientpreparationsservice "github.com/prixfixeco/api_server/internal/services/validingredientpreparations"
	validingredientsservice "github.com/prixfixeco/api_server/internal/services/validingredients"
	validinstrumentsservice "github.com/prixfixeco/api_server/internal/services/validinstruments"
	validpreparationsservice "github.com/prixfixeco/api_server/internal/services/validpreparations"
	webhooksservice "github.com/prixfixeco/api_server/internal/services/webhooks"
	"github.com/prixfixeco/api_server/internal/storage"
	"github.com/prixfixeco/api_server/internal/uploads"
	"github.com/prixfixeco/api_server/internal/uploads/images"
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
		uploads.Providers,
		storage.Providers,
		chi.Providers,
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
		mealplanoptionsservice.Providers,
		mealplanoptionvotesservice.Providers,
	)

	return nil, nil
}
