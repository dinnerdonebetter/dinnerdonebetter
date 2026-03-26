package http

import (
	paymentsmanager "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/payments/manager"

	"github.com/samber/do/v2"
	"github.com/verygoodsoftwarenotvirus/platform/v3/observability/logging"
	"github.com/verygoodsoftwarenotvirus/platform/v3/observability/tracing"
)

// RegisterPaymentsHTTP registers the payments HTTP handler with the injector.
func RegisterPaymentsHTTP(i do.Injector) {
	do.ProvideValue(i, WebhookSignatureHeader(StripeSignatureHeader))

	do.Provide[*WebhookHandler](i, func(i do.Injector) (*WebhookHandler, error) {
		return NewWebhookHandler(
			do.MustInvoke[logging.Logger](i),
			do.MustInvoke[tracing.TracerProvider](i),
			do.MustInvoke[paymentsmanager.PaymentsDataManager](i),
			do.MustInvoke[WebhookSignatureHeader](i),
		), nil
	})
}
