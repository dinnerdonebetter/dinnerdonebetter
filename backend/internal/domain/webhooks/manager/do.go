package manager

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/domain/webhooks"

	"github.com/samber/do/v2"
	"github.com/verygoodsoftwarenotvirus/platform/messagequeue"
	msgconfig "github.com/verygoodsoftwarenotvirus/platform/messagequeue/config"
	"github.com/verygoodsoftwarenotvirus/platform/observability/logging"
	"github.com/verygoodsoftwarenotvirus/platform/observability/tracing"
)

// RegisterWebhookDataManager registers the webhook data manager with the injector.
func RegisterWebhookDataManager(i do.Injector) {
	do.Provide[WebhookDataManager](i, func(i do.Injector) (WebhookDataManager, error) {
		return NewWebhookDataManager(
			do.MustInvoke[context.Context](i),
			do.MustInvoke[tracing.TracerProvider](i),
			do.MustInvoke[logging.Logger](i),
			do.MustInvoke[webhooks.Repository](i),
			do.MustInvoke[*msgconfig.QueuesConfig](i),
			do.MustInvoke[messagequeue.PublisherProvider](i),
		)
	})
}
