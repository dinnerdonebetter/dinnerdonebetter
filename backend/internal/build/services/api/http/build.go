//go:build wireinject

package api

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/authentication"
	"github.com/dinnerdonebetter/backend/internal/config"
	identitymgr "github.com/dinnerdonebetter/backend/internal/domain/identity/manager"
	paymentsmanager "github.com/dinnerdonebetter/backend/internal/domain/payments/manager"
	analyticscfg "github.com/dinnerdonebetter/backend/internal/platform/analytics/config"
	databasecfg "github.com/dinnerdonebetter/backend/internal/platform/database/config"
	"github.com/dinnerdonebetter/backend/internal/platform/encoding"
	featureflagscfg "github.com/dinnerdonebetter/backend/internal/platform/featureflags/config"
	msgconfig "github.com/dinnerdonebetter/backend/internal/platform/messagequeue/config"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	loggingcfg "github.com/dinnerdonebetter/backend/internal/platform/observability/logging/config"
	metricscfg "github.com/dinnerdonebetter/backend/internal/platform/observability/metrics/config"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	tracingcfg "github.com/dinnerdonebetter/backend/internal/platform/observability/tracing/config"
	"github.com/dinnerdonebetter/backend/internal/platform/random"
	routingcfg "github.com/dinnerdonebetter/backend/internal/platform/routing/config"
	"github.com/dinnerdonebetter/backend/internal/platform/server/http"
	"github.com/dinnerdonebetter/backend/internal/repositories"
	auditrepo "github.com/dinnerdonebetter/backend/internal/repositories/postgres/auditlogentries"
	identityrepo "github.com/dinnerdonebetter/backend/internal/repositories/postgres/identity"
	oauthrepo "github.com/dinnerdonebetter/backend/internal/repositories/postgres/oauth"
	paymentsrepo "github.com/dinnerdonebetter/backend/internal/repositories/postgres/payments"
	authservice "github.com/dinnerdonebetter/backend/internal/services/auth/handlers/authentication"
	paymentsadapters "github.com/dinnerdonebetter/backend/internal/services/payments/adapters"
	paymentshttp "github.com/dinnerdonebetter/backend/internal/services/payments/http"

	"github.com/google/wire"
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
		featureflagscfg.ProvidersFeatureFlags,
		tracing.ProvidersTracing,
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
	)

	return nil, nil
}
