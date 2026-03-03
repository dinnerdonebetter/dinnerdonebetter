package notifications

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/dinnerdonebetter/backend/internal/platform/notifications/apns"
	"github.com/dinnerdonebetter/backend/internal/platform/notifications/fcm"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
)

// ErrPlatformNotSupported is returned when attempting to send to a platform
// that has no configured sender (e.g., iOS token but APNs not configured).
var ErrPlatformNotSupported = errors.New("push notifications not configured for this platform")

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

	platform = strings.ToLower(strings.TrimSpace(platform))
	logger := s.logger.WithValue("platform", platform)

	switch platform {
	case platformIOS:
		if s.apnsSender == nil {
			return observability.PrepareAndLogError(ErrPlatformNotSupported, logger, span, "sending apns notification")
		}
		return s.apnsSender.Send(ctx, token, title, body)
	case platformAndroid:
		if s.fcmSender == nil {
			return observability.PrepareAndLogError(ErrPlatformNotSupported, logger, span, "sending apns notification")
		}
		return s.fcmSender.Send(ctx, token, title, body)
	default:
		return observability.PrepareAndLogError(fmt.Errorf("unknown platform %q", platform), logger, span, "sending apns notification")
	}
}
