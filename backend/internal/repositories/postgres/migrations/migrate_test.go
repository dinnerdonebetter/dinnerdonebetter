package migrations

import (
	"testing"
	"time"

	pgtesting "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/repositories/postgres/testing"

	"github.com/stretchr/testify/require"
	"github.com/verygoodsoftwarenotvirus/platform/v4/observability/logging"
	"github.com/verygoodsoftwarenotvirus/platform/v4/pointer"
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
