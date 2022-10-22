package recipes

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/prixfixeco/api_server/internal/uploads"
)

func TestConfig_ValidateWithContext(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		cfg := &Config{
			PublicMediaURLPrefix: t.Name(),
			Uploads:              uploads.Config{},
			DataChangesTopicName: "blah",
		}

		assert.NoError(t, cfg.ValidateWithContext(ctx))
	})
}
