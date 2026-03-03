package notifications

import (
	"context"
	"strings"

	"github.com/dinnerdonebetter/backend/internal/platform/notifications/apns"
	"github.com/dinnerdonebetter/backend/internal/platform/notifications/fcm"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
)

const (
	platformIOS     = "ios"
	platformAndroid = "android"
	o11yName        = "mobile_push_sender"
)

// MultiPlatformPushSender routes push notifications to APNs (iOS) or FCM (Android).
type MultiPlatformPushSender struct {
	tracer     tracing.Tracer
	apnsSender *apns.Sender
	fcmSender  *fcm.Sender
	logger     logging.Logger
}

// NewMultiPlatformPushSender creates a sender that routes by platform.
func NewMultiPlatformPushSender(
	apnsSender *apns.Sender,
	fcmSender *fcm.Sender,
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
) *MultiPlatformPushSender {
	return &MultiPlatformPushSender{
		apnsSender: apnsSender,
		fcmSender:  fcmSender,
		tracer:     tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(o11yName)),
		logger:     logging.EnsureLogger(logger).WithName(o11yName),
	}
}

// SendPush sends a push notification to a single device token, routing by platform.
func (s *MultiPlatformPushSender) SendPush(ctx context.Context, platform, token, title, body string) error {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithValue("platform", platform)

	platform = strings.ToLower(strings.TrimSpace(platform))

	switch platform {
	case platformIOS:
		if s.apnsSender == nil {
			logger.Debug("push: ios token but APNs sender not configured, skipping")
			return nil
		}
		return s.apnsSender.Send(ctx, token, title, body)
	case platformAndroid:
		if s.fcmSender == nil {
			logger.Debug("push: android token but FCM sender not configured, skipping")
			return nil
		}
		return s.fcmSender.Send(ctx, token, title, body)
	default:
		s.logger.WithValue("platform", platform).Debug("push: unknown platform, skipping")
		return nil
	}
}
