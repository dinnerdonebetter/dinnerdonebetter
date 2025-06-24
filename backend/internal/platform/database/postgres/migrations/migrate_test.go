package migrations

import (
	"context"
	"testing"
	"time"

	databasecfg "github.com/dinnerdonebetter/backend/internal/database/config"
	pgtesting "github.com/dinnerdonebetter/backend/internal/platform/database/postgres/testing"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQuerier_Migrate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleCreationTime := fakes.BuildFakeTime()

		exampleUser := fakes.BuildFakeUser()
		exampleUser.TwoFactorSecretVerifiedAt = nil
		exampleUser.CreatedAt = exampleCreationTime

		c, db := pgtesting.BuildDatabaseClientForTest(t)

		config := &databasecfg.Config{MaxPingAttempts: 1, PingWaitPeriod: time.Second}

		migrator := NewMigrator(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), db, config)
		migrator.migrateOnce.Do(func() {})

		ctx := context.Background()
		assert.NoError(t, migrator.Migrate(ctx))

		require.NoError(t, c.Terminate(ctx))
	})
}
