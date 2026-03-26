package api

import (
	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/auth"
	paymentswebhook "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/services/payments/http"

	"github.com/samber/do/v2"
	"github.com/verygoodsoftwarenotvirus/platform/v3/healthcheck"
	"github.com/verygoodsoftwarenotvirus/platform/v3/observability/logging"
	"github.com/verygoodsoftwarenotvirus/platform/v3/observability/metrics"
	"github.com/verygoodsoftwarenotvirus/platform/v3/observability/tracing"
	"github.com/verygoodsoftwarenotvirus/platform/v3/routing"
	routingcfg "github.com/verygoodsoftwarenotvirus/platform/v3/routing/config"
)

// RegisterAPIRouter registers the API router provider with the injector.
func RegisterAPIRouter(i do.Injector) {
	do.Provide[routing.Router](i, func(i do.Injector) (routing.Router, error) {
		return ProvideAPIRouter(
			*do.MustInvoke[*routingcfg.Config](i),
			do.MustInvoke[logging.Logger](i),
			do.MustInvoke[tracing.TracerProvider](i),
			do.MustInvoke[metrics.Provider](i),
			do.MustInvoke[auth.AuthDataService](i),
			do.MustInvoke[*paymentswebhook.WebhookHandler](i),
			do.MustInvoke[healthcheck.Registry](i),
		)
	})
}
