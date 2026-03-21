package datachangemessagehandler

import (
	"context"

	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/config"
	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/auth"
	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/dataprivacy"
	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/identity"
	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/internalops"
	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/mealplanning"
	notificationsmanager "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/notifications/manager"
	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/webhooks"
	identityindexing "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/services/identity/indexing"
	mealplanningindexing "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/services/mealplanning/indexing"

	"github.com/samber/do/v2"
	"github.com/verygoodsoftwarenotvirus/platform/analytics"
	"github.com/verygoodsoftwarenotvirus/platform/email"
	"github.com/verygoodsoftwarenotvirus/platform/encoding"
	"github.com/verygoodsoftwarenotvirus/platform/messagequeue"
	"github.com/verygoodsoftwarenotvirus/platform/notifications"
	"github.com/verygoodsoftwarenotvirus/platform/observability/logging"
	"github.com/verygoodsoftwarenotvirus/platform/observability/metrics"
	"github.com/verygoodsoftwarenotvirus/platform/observability/tracing"
	"github.com/verygoodsoftwarenotvirus/platform/uploads"
)

// RegisterAsyncDataChangeMessageHandler registers the async data change message handler with the injector.
func RegisterAsyncDataChangeMessageHandler(i do.Injector) {
	do.Provide[*AsyncDataChangeMessageHandler](i, func(i do.Injector) (*AsyncDataChangeMessageHandler, error) {
		return NewAsyncDataChangeMessageHandler(
			do.MustInvoke[context.Context](i),
			do.MustInvoke[logging.Logger](i),
			do.MustInvoke[tracing.TracerProvider](i),
			do.MustInvoke[*config.AsyncMessageHandlerConfig](i),
			do.MustInvoke[identity.Repository](i),
			do.MustInvoke[dataprivacy.Repository](i),
			do.MustInvoke[webhooks.Repository](i),
			do.MustInvoke[internalops.InternalOpsDataManager](i),
			do.MustInvoke[messagequeue.ConsumerProvider](i),
			do.MustInvoke[messagequeue.PublisherProvider](i),
			do.MustInvoke[analytics.EventReporter](i),
			do.MustInvoke[email.Emailer](i),
			do.MustInvoke[uploads.UploadManager](i),
			do.MustInvoke[metrics.Provider](i),
			do.MustInvoke[encoding.ServerEncoderDecoder](i),
			do.MustInvoke[*identityindexing.UserDataIndexer](i),
			do.MustInvoke[*mealplanningindexing.MealPlanningDataIndexer](i),
			do.MustInvoke[mealplanning.Repository](i),
			do.MustInvoke[auth.PasswordResetTokenDataManager](i),
			do.MustInvoke[notificationsmanager.NotificationsDataManager](i),
			do.MustInvoke[notifications.PushNotificationSender](i),
		)
	})
}
