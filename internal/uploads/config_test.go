package uploads

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/prixfixeco/api_server/internal/storage"
)

func TestConfig_ValidateWithContext(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		cfg := &Config{
			Storage: storage.Config{
				FilesystemConfig:  &storage.FilesystemConfig{RootDirectory: "/blah"},
				S3Config:          &storage.S3Config{BucketName: "blahs"},
				BucketName:        "blahs",
				UploadFilenameKey: "blahs",
				Provider:          "blahs",
			},
			Debug: false,
		}

		assert.NoError(t, cfg.ValidateWithContext(ctx))
	})
}
