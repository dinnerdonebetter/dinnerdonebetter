package config

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dinnerdonebetter/backend/internal/featureflags/launchdarkly"
)

func TestConfig_ValidateWithContext(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		cfg := &Config{
			LaunchDarkly: &launchdarkly.Config{
				SDKKey:      t.Name(),
				InitTimeout: 123,
			},
			Provider: ProviderLaunchDarkly,
		}

		assert.NoError(t, cfg.ValidateWithContext(ctx))
	})
}
