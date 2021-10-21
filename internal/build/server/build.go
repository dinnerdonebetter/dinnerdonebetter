//go:build wireinject
// +build wireinject

package server

import (
	"context"

	"github.com/google/wire"

	"gitlab.com/prixfixe/prixfixe/internal/authentication"
	"gitlab.com/prixfixe/prixfixe/internal/config"
	database "gitlab.com/prixfixe/prixfixe/internal/database"
	dbconfig "gitlab.com/prixfixe/prixfixe/internal/database/config"
	"gitlab.com/prixfixe/prixfixe/internal/encoding"
	msgconfig "gitlab.com/prixfixe/prixfixe/internal/messagequeue/config"
	observability "gitlab.com/prixfixe/prixfixe/internal/observability"
	"gitlab.com/prixfixe/prixfixe/internal/observability/logging"
	"gitlab.com/prixfixe/prixfixe/internal/observability/metrics"
	chi "gitlab.com/prixfixe/prixfixe/internal/routing/chi"
	"gitlab.com/prixfixe/prixfixe/internal/search/elasticsearch"
	server "gitlab.com/prixfixe/prixfixe/internal/server"
	accountsservice "gitlab.com/prixfixe/prixfixe/internal/services/accounts"
	adminservice "gitlab.com/prixfixe/prixfixe/internal/services/admin"
	apiclientsservice "gitlab.com/prixfixe/prixfixe/internal/services/apiclients"
	authservice "gitlab.com/prixfixe/prixfixe/internal/services/authentication"
	frontendservice "gitlab.com/prixfixe/prixfixe/internal/services/frontend"
	mealplanoptionsservice "gitlab.com/prixfixe/prixfixe/internal/services/mealplanoptions"
	mealplanoptionvotesservice "gitlab.com/prixfixe/prixfixe/internal/services/mealplanoptionvotes"
	mealplansservice "gitlab.com/prixfixe/prixfixe/internal/services/mealplans"
	recipesservice "gitlab.com/prixfixe/prixfixe/internal/services/recipes"
	recipestepingredientsservice "gitlab.com/prixfixe/prixfixe/internal/services/recipestepingredients"
	recipestepinstrumentsservice "gitlab.com/prixfixe/prixfixe/internal/services/recipestepinstruments"
	recipestepproductsservice "gitlab.com/prixfixe/prixfixe/internal/services/recipestepproducts"
	recipestepsservice "gitlab.com/prixfixe/prixfixe/internal/services/recipesteps"
	usersservice "gitlab.com/prixfixe/prixfixe/internal/services/users"
	validingredientpreparationsservice "gitlab.com/prixfixe/prixfixe/internal/services/validingredientpreparations"
	validingredientsservice "gitlab.com/prixfixe/prixfixe/internal/services/validingredients"
	validinstrumentsservice "gitlab.com/prixfixe/prixfixe/internal/services/validinstruments"
	validpreparationsservice "gitlab.com/prixfixe/prixfixe/internal/services/validpreparations"
	webhooksservice "gitlab.com/prixfixe/prixfixe/internal/services/webhooks"
	websocketsservice "gitlab.com/prixfixe/prixfixe/internal/services/websockets"
	storage "gitlab.com/prixfixe/prixfixe/internal/storage"
	uploads "gitlab.com/prixfixe/prixfixe/internal/uploads"
	images "gitlab.com/prixfixe/prixfixe/internal/uploads/images"
)

// Build builds a server.
func Build(
	ctx context.Context,
	logger logging.Logger,
	cfg *config.InstanceConfig,
) (*server.HTTPServer, error) {
	wire.Build(
		elasticsearch.Providers,
		config.Providers,
		database.Providers,
		dbconfig.Providers,
		encoding.Providers,
		msgconfig.Providers,
		server.Providers,
		metrics.Providers,
		images.Providers,
		uploads.Providers,
		observability.Providers,
		storage.Providers,
		chi.Providers,
		authentication.Providers,
		authservice.Providers,
		usersservice.Providers,
		accountsservice.Providers,
		apiclientsservice.Providers,
		webhooksservice.Providers,
		websocketsservice.Providers,
		adminservice.Providers,
		frontendservice.Providers,
		validinstrumentsservice.Providers,
		validingredientsservice.Providers,
		validpreparationsservice.Providers,
		validingredientpreparationsservice.Providers,
		recipesservice.Providers,
		recipestepsservice.Providers,
		recipestepinstrumentsservice.Providers,
		recipestepingredientsservice.Providers,
		recipestepproductsservice.Providers,
		mealplansservice.Providers,
		mealplanoptionsservice.Providers,
		mealplanoptionvotesservice.Providers,
	)

	return nil, nil
}
