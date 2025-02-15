package config

import (
	"context"
	"testing"

	msgconfig "github.com/dinnerdonebetter/backend/internal/lib/messagequeue/config"

	"github.com/stretchr/testify/assert"
)

func TestQueueSettings_ValidateWithContext(T *testing.T) {
	T.Parallel()

	T.Run("invalid", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		cfg := msgconfig.QueuesConfig{}

		assert.Error(t, cfg.ValidateWithContext(ctx))
	})
}
