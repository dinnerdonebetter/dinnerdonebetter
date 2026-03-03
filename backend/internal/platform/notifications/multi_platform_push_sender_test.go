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

	t.Run("noop when both senders nil", func(t *testing.T) {
		t.Parallel()

		sender := NewMultiPlatformPushSender(nil, nil, logger, tracer)
		err := sender.SendPush(ctx, "ios", "token", "title", "body")
		assert.NoError(t, err)
	})

	t.Run("unknown platform returns nil", func(t *testing.T) {
		t.Parallel()

		sender := NewMultiPlatformPushSender(nil, nil, logger, tracer)
		err := sender.SendPush(ctx, "unknown", "token", "title", "body")
		assert.NoError(t, err)
	})
}
