//go:build wireinject

package api

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/authentication"
	"github.com/dinnerdonebetter/backend/internal/config"
	identitymgr "github.com/dinnerdonebetter/backend/internal/domain/identity/manager"
	paymentsmanager "github.com/dinnerdonebetter/backend/internal/domain/payments/manager"
	"github.com/dinnerdonebetter/backend/internal/repositories"
	auditrepo "github.com/dinnerdonebetter/backend/internal/repositories/postgres/auditlogentries"
	identityrepo "github.com/dinnerdonebetter/backend/internal/repositories/postgres/identity"
	oauthrepo "github.com/dinnerdonebetter/backend/internal/repositories/postgres/oauth"
	paymentsrepo "github.com/dinnerdonebetter/backend/internal/repositories/postgres/payments"
	authservice "github.com/dinnerdonebetter/backend/internal/services/auth/handlers/authentication"
	paymentsadapters "github.com/dinnerdonebetter/backend/internal/services/payments/adapters"
	paymentshttp "github.com/dinnerdonebetter/backend/internal/services/payments/http"

	"github.com/google/wire"
	analyticscfg "github.com/verygoodsoftwarenotvirus/platform/analytics/config"
	databasecfg "github.com/verygoodsoftwarenotvirus/platform/database/config"
	"github.com/verygoodsoftwarenotvirus/platform/encoding"
	msgconfig "github.com/verygoodsoftwarenotvirus/platform/messagequeue/config"
	"github.com/verygoodsoftwarenotvirus/platform/observability"
	loggingcfg "github.com/verygoodsoftwarenotvirus/platform/observability/logging/config"
	metricscfg "github.com/verygoodsoftwarenotvirus/platform/observability/metrics/config"
	tracingcfg "github.com/verygoodsoftwarenotvirus/platform/observability/tracing/config"
	"github.com/verygoodsoftwarenotvirus/platform/random"
	routingcfg "github.com/verygoodsoftwarenotvirus/platform/routing/config"
	"github.com/verygoodsoftwarenotvirus/platform/server/http"
)

// Build builds a server.
func Build(
	ctx context.Context,
	cfg *config.APIServiceConfig,
) (http.Server, error) {
	wire.Build(
		authentication.AuthProviders,
		encoding.Providers,
		msgconfig.MessageQueueProviders,
		analyticscfg.Providers,
		tracingcfg.TracingConfigProviders,
		observability.O11yProviders,
		databasecfg.DatabaseConfigProviders,
		repositories.RepositoryProviders,
		loggingcfg.LogConfigProviders,
		authservice.AuthHTTPServiceProviders,
		metricscfg.MetricsConfigProviders,
		http.ProvidersHTTP,
		routingcfg.RoutingConfigProviders,
		// repos
		auditrepo.AuditRepoProviders,
		identityrepo.IDRepoProviders,
		oauthrepo.OAuthRepoProviders,
		paymentsrepo.PaymentsRepoProviders,
		// manager
		random.RandProviders,
		identitymgr.IDManagerProviders,
		paymentsmanager.PaymentsManagerProviders,
		paymentsadapters.PaymentsAdapterProviders,
		// payments http
		paymentshttp.PaymentsHTTPProviders,
		ProvideTextSearchConfig,
		ProvideUserTextSearcher,
		ConfigProviders,
		ProvideAPIRouter,
		wire.Value("api_server"), // HTTP server logger service name
	)

	return nil, nil
}
