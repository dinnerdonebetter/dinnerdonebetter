//go:build wireinject

package outboundemailhandler

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/config"
	"github.com/dinnerdonebetter/backend/internal/functions/outboundemailhandler"
	analyticscfg "github.com/dinnerdonebetter/backend/internal/platform/analytics/config"
	emailcfg "github.com/dinnerdonebetter/backend/internal/platform/email/config"
	msgconfig "github.com/dinnerdonebetter/backend/internal/platform/messagequeue/config"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	loggingcfg "github.com/dinnerdonebetter/backend/internal/platform/observability/logging/config"
	metricscfg "github.com/dinnerdonebetter/backend/internal/platform/observability/metrics/config"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	tracingcfg "github.com/dinnerdonebetter/backend/internal/platform/observability/tracing/config"

	"github.com/google/wire"
)

// Build builds the outbound email handler.
func Build(
	ctx context.Context,
	cfg *config.OutboundEmailHandlerConfig,
) (*outboundemailhandler.OutboundEmailHandler, error) {
	wire.Build(
		outboundemailhandler.Providers,
		msgconfig.ProvideConsumerProvider,
		analyticscfg.Providers,
		emailcfg.Providers,
		metricscfg.MetricsConfigProviders,
		loggingcfg.LogConfigProviders,
		tracingcfg.TracingConfigProviders,
		observability.O11yProviders,
		tracing.ProvidersTracing,
		ConfigProviders,
	)

	return nil, nil
}
