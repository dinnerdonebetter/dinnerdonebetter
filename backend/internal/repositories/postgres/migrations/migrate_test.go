package migrations

import (
	"testing"
	"time"

	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/pointer"
	pgtesting "github.com/dinnerdonebetter/backend/internal/repositories/postgres/testing"

	"github.com/stretchr/testify/require"
)

func TestQuerier_Migrate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		container, db, _ := pgtesting.BuildDatabaseContainerForTest(t)
		require.NoError(t, NewMigrator(logging.NewNoopLogger()).Migrate(ctx, db))

		if err := container.Stop(ctx, pointer.To(time.Second*10)); err != nil {
			t.Log(err)
		}
	})
}
