//go:build wireinject

package emaildeliverabilitytest

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/config"
	emailcfg "github.com/dinnerdonebetter/backend/internal/platform/email/config"
	"github.com/dinnerdonebetter/backend/internal/platform/httpclient"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	loggingcfg "github.com/dinnerdonebetter/backend/internal/platform/observability/logging/config"
	metricscfg "github.com/dinnerdonebetter/backend/internal/platform/observability/metrics/config"
	tracingcfg "github.com/dinnerdonebetter/backend/internal/platform/observability/tracing/config"
	emaildeliverabilitytest "github.com/dinnerdonebetter/backend/internal/services/email/workers/email_deliverability_test"

	"github.com/google/wire"
)

// Build builds the email deliverability test job.
func Build(
	ctx context.Context,
	cfg *config.EmailDeliverabilityTestConfig,
) (*emaildeliverabilitytest.Job, error) {
	wire.Build(
		emaildeliverabilitytest.ProvidersEmailDeliverabilityTest,
		tracingcfg.TracingConfigProviders,
		observability.O11yProviders,
		loggingcfg.LogConfigProviders,
		metricscfg.MetricsConfigProviders,
		httpclient.Providers,
		emailcfg.Providers,
		ConfigProviders,
	)

	return nil, nil
}
