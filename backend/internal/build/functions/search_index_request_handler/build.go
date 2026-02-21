//go:build wireinject

package searchindexrequesthandler

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/config"
	"github.com/dinnerdonebetter/backend/internal/functions/searchindexrequesthandler"
	databasecfg "github.com/dinnerdonebetter/backend/internal/platform/database/config"
	"github.com/dinnerdonebetter/backend/internal/platform/database/postgres"
	msgconfig "github.com/dinnerdonebetter/backend/internal/platform/messagequeue/config"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	loggingcfg "github.com/dinnerdonebetter/backend/internal/platform/observability/logging/config"
	metricscfg "github.com/dinnerdonebetter/backend/internal/platform/observability/metrics/config"
	tracingcfg "github.com/dinnerdonebetter/backend/internal/platform/observability/tracing/config"
	"github.com/dinnerdonebetter/backend/internal/repositories/postgres/auditlogentries"
	"github.com/dinnerdonebetter/backend/internal/repositories/postgres/identity"
	"github.com/dinnerdonebetter/backend/internal/repositories/postgres/mealplanning"
	identityindexing "github.com/dinnerdonebetter/backend/internal/services/identity/indexing"
	eatingindexing "github.com/dinnerdonebetter/backend/internal/services/mealplanning/indexing"

	"github.com/google/wire"
)

// Build builds the search index request handler.
func Build(
	ctx context.Context,
	cfg *config.SearchIndexRequestHandlerConfig,
) (*searchindexrequesthandler.SearchIndexRequestHandler, error) {
	wire.Build(
		searchindexrequesthandler.Providers,
		msgconfig.ProvideConsumerProvider,
		databasecfg.ClientConfigProviders,
		postgres.PGProviders,
		auditlogentries.AuditRepoProviders,
		identity.IDRepoProviders,
		mealplanning.MPRepoProviders,
		metricscfg.MetricsConfigProviders,
		loggingcfg.LogConfigProviders,
		tracingcfg.TracingConfigProviders,
		observability.O11yProviders,
		identityindexing.Providers,
		eatingindexing.Providers,
		ConfigProviders,
		SearcherProviders,
	)

	return nil, nil
}
