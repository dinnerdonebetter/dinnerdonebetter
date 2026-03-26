package payments

import (
	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/audit"
	domainpayments "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/payments"

	"github.com/samber/do/v2"
	"github.com/verygoodsoftwarenotvirus/platform/v3/database"
	"github.com/verygoodsoftwarenotvirus/platform/v3/observability/logging"
	"github.com/verygoodsoftwarenotvirus/platform/v3/observability/tracing"
)

// RegisterPaymentsRepository registers the payments repository with the injector.
func RegisterPaymentsRepository(i do.Injector) {
	do.Provide[domainpayments.Repository](i, func(i do.Injector) (domainpayments.Repository, error) {
		return ProvidePaymentsRepository(
			do.MustInvoke[logging.Logger](i),
			do.MustInvoke[tracing.TracerProvider](i),
			do.MustInvoke[audit.Repository](i),
			do.MustInvoke[database.Client](i),
		), nil
	})
}
