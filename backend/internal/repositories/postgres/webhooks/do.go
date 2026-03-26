package webhooks

import (
	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/audit"
	domainwebhooks "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/webhooks"

	"github.com/samber/do/v2"
	"github.com/verygoodsoftwarenotvirus/platform/v3/database"
	"github.com/verygoodsoftwarenotvirus/platform/v3/observability/logging"
	"github.com/verygoodsoftwarenotvirus/platform/v3/observability/tracing"
)

// RegisterWebhooksRepository registers the webhooks repository with the injector.
func RegisterWebhooksRepository(i do.Injector) {
	do.Provide[domainwebhooks.Repository](i, func(i do.Injector) (domainwebhooks.Repository, error) {
		return ProvideWebhooksRepository(
			do.MustInvoke[logging.Logger](i),
			do.MustInvoke[tracing.TracerProvider](i),
			do.MustInvoke[audit.Repository](i),
			do.MustInvoke[database.Client](i),
		), nil
	})
}
