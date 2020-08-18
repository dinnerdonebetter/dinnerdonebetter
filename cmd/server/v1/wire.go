// +build wireinject

package main

import (
	"context"
	"database/sql"

	database "gitlab.com/prixfixe/prixfixe/database/v1"
	"gitlab.com/prixfixe/prixfixe/internal/v1/auth"
	"gitlab.com/prixfixe/prixfixe/internal/v1/config"
	"gitlab.com/prixfixe/prixfixe/internal/v1/encoding"
	"gitlab.com/prixfixe/prixfixe/internal/v1/metrics"
	server "gitlab.com/prixfixe/prixfixe/server/v1"
	httpserver "gitlab.com/prixfixe/prixfixe/server/v1/http"
	authservice "gitlab.com/prixfixe/prixfixe/services/v1/auth"
	frontendservice "gitlab.com/prixfixe/prixfixe/services/v1/frontend"
	invitationsservice "gitlab.com/prixfixe/prixfixe/services/v1/invitations"
	iterationmediasservice "gitlab.com/prixfixe/prixfixe/services/v1/iterationmedias"
	oauth2clientsservice "gitlab.com/prixfixe/prixfixe/services/v1/oauth2clients"
	recipeiterationsservice "gitlab.com/prixfixe/prixfixe/services/v1/recipeiterations"
	recipesservice "gitlab.com/prixfixe/prixfixe/services/v1/recipes"
	recipestepeventsservice "gitlab.com/prixfixe/prixfixe/services/v1/recipestepevents"
	recipestepingredientsservice "gitlab.com/prixfixe/prixfixe/services/v1/recipestepingredients"
	recipestepinstrumentsservice "gitlab.com/prixfixe/prixfixe/services/v1/recipestepinstruments"
	recipestepproductsservice "gitlab.com/prixfixe/prixfixe/services/v1/recipestepproducts"
	recipestepsservice "gitlab.com/prixfixe/prixfixe/services/v1/recipesteps"
	reportsservice "gitlab.com/prixfixe/prixfixe/services/v1/reports"
	requiredpreparationinstrumentsservice "gitlab.com/prixfixe/prixfixe/services/v1/requiredpreparationinstruments"
	usersservice "gitlab.com/prixfixe/prixfixe/services/v1/users"
	validingredientpreparationsservice "gitlab.com/prixfixe/prixfixe/services/v1/validingredientpreparations"
	validingredientsservice "gitlab.com/prixfixe/prixfixe/services/v1/validingredients"
	validinstrumentsservice "gitlab.com/prixfixe/prixfixe/services/v1/validinstruments"
	validpreparationsservice "gitlab.com/prixfixe/prixfixe/services/v1/validpreparations"
	webhooksservice "gitlab.com/prixfixe/prixfixe/services/v1/webhooks"

	"github.com/google/wire"
	"gitlab.com/verygoodsoftwarenotvirus/logging/v1"
	"gitlab.com/verygoodsoftwarenotvirus/newsman"
)

// ProvideReporter is an obligatory function that hopefully wire will eliminate for me one day.
func ProvideReporter(n *newsman.Newsman) newsman.Reporter {
	return n
}

// BuildServer builds a server.
func BuildServer(
	ctx context.Context,
	cfg *config.ServerConfig,
	logger logging.Logger,
	database database.DataManager,
	db *sql.DB,
) (*server.Server, error) {
	wire.Build(
		config.Providers,
		auth.Providers,
		// server things,
		server.Providers,
		encoding.Providers,
		httpserver.Providers,
		// metrics,
		metrics.Providers,
		// external libs,
		newsman.NewNewsman,
		ProvideReporter,
		// services,
		authservice.Providers,
		usersservice.Providers,
		validinstrumentsservice.Providers,
		validingredientsservice.Providers,
		validpreparationsservice.Providers,
		validingredientpreparationsservice.Providers,
		requiredpreparationinstrumentsservice.Providers,
		recipesservice.Providers,
		recipestepsservice.Providers,
		recipestepinstrumentsservice.Providers,
		recipestepingredientsservice.Providers,
		recipestepproductsservice.Providers,
		recipeiterationsservice.Providers,
		recipestepeventsservice.Providers,
		iterationmediasservice.Providers,
		invitationsservice.Providers,
		reportsservice.Providers,
		frontendservice.Providers,
		webhooksservice.Providers,
		oauth2clientsservice.Providers,
	)
	return nil, nil
}
