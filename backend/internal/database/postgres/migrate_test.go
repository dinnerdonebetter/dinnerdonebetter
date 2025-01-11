package postgres

import (
	"context"
	"testing"
	"time"

	databasecfg "github.com/dinnerdonebetter/backend/internal/database/config"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestQuerier_Migrate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleCreationTime := fakes.BuildFakeTime()

		exampleUser := fakes.BuildFakeUser()
		exampleUser.TwoFactorSecretVerifiedAt = nil
		exampleUser.CreatedAt = exampleCreationTime

		ctx := context.Background()
		c, db := buildTestClient(t)

		c.timeFunc = func() time.Time {
			return exampleCreationTime
		}
		c.config = &databasecfg.Config{MaxPingAttempts: 1, PingWaitPeriod: time.Second}

		// called by c.IsReady()
		db.ExpectPing()

		c.migrateOnce.Do(func() {})

		err := c.Migrate(ctx)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, db)
	})
}
