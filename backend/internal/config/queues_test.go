package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	msgconfig "github.com/verygoodsoftwarenotvirus/platform/v3/messagequeue/config"
)

func TestQueueSettings_ValidateWithContext(T *testing.T) {
	T.Parallel()

	T.Run("invalid", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		cfg := msgconfig.QueuesConfig{}

		assert.Error(t, cfg.ValidateWithContext(ctx))
	})
}
