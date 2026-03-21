//go:build wireinject

package emaildeliverabilitytest

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/config"
	emaildeliverabilitytest "github.com/dinnerdonebetter/backend/internal/services/email/workers/email_deliverability_test"

	"github.com/google/wire"
	emailcfg "github.com/verygoodsoftwarenotvirus/platform/email/config"
	"github.com/verygoodsoftwarenotvirus/platform/httpclient"
	"github.com/verygoodsoftwarenotvirus/platform/observability"
	loggingcfg "github.com/verygoodsoftwarenotvirus/platform/observability/logging/config"
	metricscfg "github.com/verygoodsoftwarenotvirus/platform/observability/metrics/config"
	tracingcfg "github.com/verygoodsoftwarenotvirus/platform/observability/tracing/config"
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
