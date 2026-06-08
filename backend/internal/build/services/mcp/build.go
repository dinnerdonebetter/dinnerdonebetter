package mcpbuild

import (
	"context"

	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/authentication"
	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/config"
	mealplanningregistration "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/mealplanning/registration"
	auditrepo "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/repositories/postgres/auditlogentries"
	identityrepo "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/repositories/postgres/identity"
	issuereportsrepo "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/repositories/postgres/issuereports"
	waitlistsrepo "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/repositories/postgres/waitlists"
	webhooksrepo "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/repositories/postgres/webhooks"

	databasecfg "github.com/primandproper/platform/database/config"
	"github.com/primandproper/platform/database/postgres"
	"github.com/primandproper/platform/observability"
	loggingcfg "github.com/primandproper/platform/observability/logging/config"
	metricscfg "github.com/primandproper/platform/observability/metrics/config"
	tracingcfg "github.com/primandproper/platform/observability/tracing/config"

	"github.com/samber/do/v2"
)

// BuildInjector creates and configures the dependency injection container for the MCP server.
func BuildInjector(ctx context.Context, cfg *config.MCPServiceConfig) *do.RootScope {
	i := do.New()

	do.ProvideValue(i, ctx)
	do.ProvideValue(i, cfg)

	// config field extraction
	RegisterConfigs(i)

	// platform providers
	observability.RegisterO11yConfigs(i)
	metricscfg.RegisterMetricsProvider(i)
	loggingcfg.RegisterLogger(i)
	tracingcfg.RegisterTracerProvider(i)
	databasecfg.RegisterClientConfig(i)
	postgres.RegisterDatabaseClient(i)

	// authentication (for login credential validation)
	authentication.RegisterAuth(i)

	// repositories
	auditrepo.RegisterAuditLogRepository(i)
	identityrepo.RegisterIdentityRepository(i)
	webhooksrepo.RegisterWebhooksRepository(i)
	waitlistsrepo.RegisterWaitlistsRepository(i)
	issuereportsrepo.RegisterIssueReportsRepository(i)

	// Domain: mealplanning
	mealplanningregistration.RegisterForMCPServer(i)

	return i
}
