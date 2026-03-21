package api

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/authentication"
	authcfg "github.com/dinnerdonebetter/backend/internal/authentication/config"
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

	"github.com/samber/do/v2"
	analyticscfg "github.com/verygoodsoftwarenotvirus/platform/analytics/config"
	"github.com/verygoodsoftwarenotvirus/platform/database/postgres"
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
	i := do.New()

	do.ProvideValue(i, ctx)
	do.ProvideValue(i, cfg)

	// config field extraction
	RegisterConfigs(i)

	// platform providers
	observability.RegisterO11yConfigs(i)
	loggingcfg.RegisterLogger(i)
	tracingcfg.RegisterTracerProvider(i)
	metricscfg.RegisterMetricsProvider(i)
	encoding.RegisterServerEncoderDecoder(i)
	msgconfig.RegisterMessageQueue(i)
	analyticscfg.RegisterEventReporter(i)
	postgres.RegisterDatabaseClient(i)
	routingcfg.RegisterRouteParamManager(i)
	random.RegisterGenerator(i)
	http.RegisterHTTPServer(i, "api_server")

	// authentication
	authentication.RegisterAuth(i)
	authcfg.RegisterConfigs(i)

	// repos
	repositories.RegisterMigrator(i)
	auditrepo.RegisterAuditLogRepository(i)
	identityrepo.RegisterIdentityRepository(i)
	oauthrepo.RegisterOAuthRepository(i)
	paymentsrepo.RegisterPaymentsRepository(i)

	// managers
	identitymgr.RegisterIdentityDataManager(i)
	paymentsmanager.RegisterPaymentsDataManager(i)
	paymentsadapters.RegisterPaymentProcessorRegistry(i)

	// services
	authservice.RegisterAuthHTTPService(i)
	paymentshttp.RegisterPaymentsHTTP(i)

	// searchers & routes
	RegisterSearchers(i)
	RegisterAPIRouter(i)

	return do.MustInvoke[http.Server](i), nil
}
