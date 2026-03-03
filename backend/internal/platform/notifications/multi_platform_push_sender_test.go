package notifications

import (
	"context"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"

	"github.com/stretchr/testify/assert"
)

func TestMultiPlatformPushSender_SendPush(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	logger := logging.NewNoopLogger()
	tracer := tracing.NewNoopTracerProvider()

	t.Run("ios returns ErrPlatformNotSupported when apnsSender nil", func(t *testing.T) {
		t.Parallel()

		sender := NewMultiPlatformPushSender(nil, nil, logger, tracer)
		err := sender.SendPush(ctx, "ios", "token", "title", "body")
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrPlatformNotSupported)
	})

	t.Run("android returns ErrPlatformNotSupported when fcmSender nil", func(t *testing.T) {
		t.Parallel()

		sender := NewMultiPlatformPushSender(nil, nil, logger, tracer)
		err := sender.SendPush(ctx, "android", "token", "title", "body")
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrPlatformNotSupported)
	})

	t.Run("unknown platform returns error", func(t *testing.T) {
		t.Parallel()

		sender := NewMultiPlatformPushSender(nil, nil, logger, tracer)
		err := sender.SendPush(ctx, "unknown", "token", "title", "body")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unknown platform")
	})
}
