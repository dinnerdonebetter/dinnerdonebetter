//go:build wireinject
// +build wireinject

package server

import (
	"context"

	"gitlab.com/prixfixe/prixfixe/internal/authentication"
	"gitlab.com/prixfixe/prixfixe/internal/config"
	database "gitlab.com/prixfixe/prixfixe/internal/database"
	dbconfig "gitlab.com/prixfixe/prixfixe/internal/database/config"
	"gitlab.com/prixfixe/prixfixe/internal/encoding"
	observability "gitlab.com/prixfixe/prixfixe/internal/observability"
	"gitlab.com/prixfixe/prixfixe/internal/observability/logging"
	"gitlab.com/prixfixe/prixfixe/internal/observability/metrics"
	chi "gitlab.com/prixfixe/prixfixe/internal/routing/chi"
	"gitlab.com/prixfixe/prixfixe/internal/search/bleve"
	server "gitlab.com/prixfixe/prixfixe/internal/server"
	adminservice "gitlab.com/prixfixe/prixfixe/internal/services/admin"
	apiclientsservice "gitlab.com/prixfixe/prixfixe/internal/services/apiclients"
	auditservice "gitlab.com/prixfixe/prixfixe/internal/services/audit"
	authservice "gitlab.com/prixfixe/prixfixe/internal/services/authentication"
	frontendservice "gitlab.com/prixfixe/prixfixe/internal/services/frontend"
	householdsservice "gitlab.com/prixfixe/prixfixe/internal/services/households"
	invitationsservice "gitlab.com/prixfixe/prixfixe/internal/services/invitations"
	recipesservice "gitlab.com/prixfixe/prixfixe/internal/services/recipes"
	recipestepingredientsservice "gitlab.com/prixfixe/prixfixe/internal/services/recipestepingredients"
	recipestepproductsservice "gitlab.com/prixfixe/prixfixe/internal/services/recipestepproducts"
	recipestepsservice "gitlab.com/prixfixe/prixfixe/internal/services/recipesteps"
	reportsservice "gitlab.com/prixfixe/prixfixe/internal/services/reports"
	usersservice "gitlab.com/prixfixe/prixfixe/internal/services/users"
	validingredientpreparationsservice "gitlab.com/prixfixe/prixfixe/internal/services/validingredientpreparations"
	validingredientsservice "gitlab.com/prixfixe/prixfixe/internal/services/validingredients"
	validinstrumentsservice "gitlab.com/prixfixe/prixfixe/internal/services/validinstruments"
	validpreparationinstrumentsservice "gitlab.com/prixfixe/prixfixe/internal/services/validpreparationinstruments"
	validpreparationsservice "gitlab.com/prixfixe/prixfixe/internal/services/validpreparations"
	webhooksservice "gitlab.com/prixfixe/prixfixe/internal/services/webhooks"
	storage "gitlab.com/prixfixe/prixfixe/internal/storage"
	uploads "gitlab.com/prixfixe/prixfixe/internal/uploads"
	images "gitlab.com/prixfixe/prixfixe/internal/uploads/images"

	"github.com/google/wire"
)

// Build builds a server.
func Build(
	ctx context.Context,
	cfg *config.InstanceConfig,
	logger logging.Logger,
) (*server.HTTPServer, error) {
	wire.Build(
		bleve.Providers,
		config.Providers,
		database.Providers,
		dbconfig.Providers,
		encoding.Providers,
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
		householdsservice.Providers,
		apiclientsservice.Providers,
		webhooksservice.Providers,
		auditservice.Providers,
		adminservice.Providers,
		frontendservice.Providers,
		validinstrumentsservice.Providers,
		validpreparationsservice.Providers,
		validingredientsservice.Providers,
		validingredientpreparationsservice.Providers,
		validpreparationinstrumentsservice.Providers,
		recipesservice.Providers,
		recipestepsservice.Providers,
		recipestepingredientsservice.Providers,
		recipestepproductsservice.Providers,
		invitationsservice.Providers,
		reportsservice.Providers,
	)

	return nil, nil
}
