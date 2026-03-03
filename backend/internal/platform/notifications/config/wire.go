package config

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/platform/notifications"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"

	"github.com/google/wire"
)

var (
	// Providers are what we provide to dependency injection.
	Providers = wire.NewSet(
		ProvidePushSender,
	)
)

// ProvidePushSender provides a PushNotificationSender from config.
func ProvidePushSender(
	ctx context.Context,
	cfg Config,
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
) (notifications.PushNotificationSender, error) {
	return (&cfg).ProvidePushSender(ctx, logger, tracerProvider)
}
