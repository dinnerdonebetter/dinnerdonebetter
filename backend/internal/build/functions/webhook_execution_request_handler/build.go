//go:build wireinject

package webhookexecutionrequesthandler

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/config"
	"github.com/dinnerdonebetter/backend/internal/functions/webhookexecutionrequesthandler"
	databasecfg "github.com/dinnerdonebetter/backend/internal/platform/database/config"
	"github.com/dinnerdonebetter/backend/internal/platform/database/postgres"
	"github.com/dinnerdonebetter/backend/internal/platform/encoding"
	msgconfig "github.com/dinnerdonebetter/backend/internal/platform/messagequeue/config"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	loggingcfg "github.com/dinnerdonebetter/backend/internal/platform/observability/logging/config"
	metricscfg "github.com/dinnerdonebetter/backend/internal/platform/observability/metrics/config"
	tracingcfg "github.com/dinnerdonebetter/backend/internal/platform/observability/tracing/config"
	"github.com/dinnerdonebetter/backend/internal/repositories/postgres/auditlogentries"
	"github.com/dinnerdonebetter/backend/internal/repositories/postgres/identity"
	"github.com/dinnerdonebetter/backend/internal/repositories/postgres/webhooks"

	"github.com/google/wire"
)

// Build builds the webhook execution request handler.
func Build(
	ctx context.Context,
	cfg *config.WebhookExecutionRequestHandlerConfig,
) (*webhookexecutionrequesthandler.WebhookExecutionRequestHandler, error) {
	wire.Build(
		webhookexecutionrequesthandler.Providers,
		msgconfig.ProvideConsumerProvider,
		databasecfg.ClientConfigProviders,
		postgres.PGProviders,
		auditlogentries.AuditRepoProviders,
		identity.IDRepoProviders,
		webhooks.WebhookProviders,
		encoding.Providers,
		metricscfg.MetricsConfigProviders,
		loggingcfg.LogConfigProviders,
		tracingcfg.TracingConfigProviders,
		observability.O11yProviders,
		ConfigProviders,
	)

	return nil, nil
}
