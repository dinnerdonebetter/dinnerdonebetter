// +build wireinject

package main

import (
	"context"

	database "gitlab.com/prixfixe/prixfixe/database/v1"
	"gitlab.com/prixfixe/prixfixe/internal/v1/auth"
	"gitlab.com/prixfixe/prixfixe/internal/v1/config"
	"gitlab.com/prixfixe/prixfixe/internal/v1/encoding"
	"gitlab.com/prixfixe/prixfixe/internal/v1/metrics"
	server "gitlab.com/prixfixe/prixfixe/server/v1"
	httpserver "gitlab.com/prixfixe/prixfixe/server/v1/http"
	auth1 "gitlab.com/prixfixe/prixfixe/services/v1/auth"
	"gitlab.com/prixfixe/prixfixe/services/v1/frontend"
	"gitlab.com/prixfixe/prixfixe/services/v1/ingredients"
	"gitlab.com/prixfixe/prixfixe/services/v1/instruments"
	"gitlab.com/prixfixe/prixfixe/services/v1/invitations"
	"gitlab.com/prixfixe/prixfixe/services/v1/iterationmedias"
	"gitlab.com/prixfixe/prixfixe/services/v1/oauth2clients"
	"gitlab.com/prixfixe/prixfixe/services/v1/preparations"
	"gitlab.com/prixfixe/prixfixe/services/v1/recipeiterations"
	"gitlab.com/prixfixe/prixfixe/services/v1/recipes"
	"gitlab.com/prixfixe/prixfixe/services/v1/recipestepevents"
	"gitlab.com/prixfixe/prixfixe/services/v1/recipestepingredients"
	"gitlab.com/prixfixe/prixfixe/services/v1/recipestepinstruments"
	"gitlab.com/prixfixe/prixfixe/services/v1/recipestepproducts"
	"gitlab.com/prixfixe/prixfixe/services/v1/recipesteps"
	"gitlab.com/prixfixe/prixfixe/services/v1/reports"
	"gitlab.com/prixfixe/prixfixe/services/v1/requiredpreparationinstruments"
	"gitlab.com/prixfixe/prixfixe/services/v1/users"
	"gitlab.com/prixfixe/prixfixe/services/v1/webhooks"

	"github.com/google/wire"
	"gitlab.com/verygoodsoftwarenotvirus/logging/v1"
	"gitlab.com/verygoodsoftwarenotvirus/newsman"
)

// ProvideReporter is an obligatory function that hopefully wire will eliminate for me one day
func ProvideReporter(n *newsman.Newsman) newsman.Reporter {
	return n
}

// BuildServer builds a server
func BuildServer(
	ctx context.Context,
	cfg *config.ServerConfig,
	logger logging.Logger,
	database database.Database,
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
		auth1.Providers,
		users.Providers,
		instruments.Providers,
		ingredients.Providers,
		preparations.Providers,
		requiredpreparationinstruments.Providers,
		recipes.Providers,
		recipesteps.Providers,
		recipestepinstruments.Providers,
		recipestepingredients.Providers,
		recipestepproducts.Providers,
		recipeiterations.Providers,
		recipestepevents.Providers,
		iterationmedias.Providers,
		invitations.Providers,
		reports.Providers,
		frontend.Providers,
		webhooks.Providers,
		oauth2clients.Providers,
	)
	return nil, nil
}
