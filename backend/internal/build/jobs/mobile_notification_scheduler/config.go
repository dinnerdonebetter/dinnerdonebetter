package mobilenotificationscheduler

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/config"
	"github.com/dinnerdonebetter/backend/internal/platform/messagequeue"
	msgconfig "github.com/dinnerdonebetter/backend/internal/platform/messagequeue/config"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"

	"github.com/google/wire"
)

var (
	// ConfigProviders represents this package's offering to the dependency injector.
	ConfigProviders = wire.NewSet(
		wire.FieldsOf(
			new(*config.MobileNotificationSchedulerConfig),
			"Queues",
			"Database",
		),
		ProvideObservabilityConfig,
		ProvideEventsConfig,
		ProvideMobileNotificationsPublisher,
		NewScheduler,
	)
)

// ProvideObservabilityConfig provides the observability config for wire.
func ProvideObservabilityConfig(cfg *config.MobileNotificationSchedulerConfig) *observability.Config {
	return &cfg.Observability
}

// ProvideEventsConfig provides the message queue events config for wire.
func ProvideEventsConfig(cfg *config.MobileNotificationSchedulerConfig) *msgconfig.Config {
	return &cfg.Events
}

// ProvideMobileNotificationsPublisher provides a publisher for the mobile_notifications topic.
func ProvideMobileNotificationsPublisher(
	ctx context.Context,
	messageQueuePublisherProvider messagequeue.PublisherProvider,
	queues *msgconfig.QueuesConfig,
) (messagequeue.Publisher, error) {
	return messageQueuePublisherProvider.ProvidePublisher(ctx, queues.MobileNotificationsTopicName)
}
