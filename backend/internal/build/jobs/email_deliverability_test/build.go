package emaildeliverabilitytest

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/config"
	emaildeliverabilitytest "github.com/dinnerdonebetter/backend/internal/services/email/workers/email_deliverability_test"

	"github.com/samber/do/v2"
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
	i := do.New()

	do.ProvideValue(i, ctx)
	do.ProvideValue(i, cfg)

	RegisterConfigs(i)

	observability.RegisterO11yConfigs(i)
	tracingcfg.RegisterTracerProvider(i)
	loggingcfg.RegisterLogger(i)
	metricscfg.RegisterMetricsProvider(i)
	httpclient.RegisterHTTPClient(i)
	emailcfg.RegisterEmailer(i)
	emaildeliverabilitytest.RegisterEmailDeliverabilityTest(i)

	return do.MustInvoke[*emaildeliverabilitytest.Job](i), nil
}
