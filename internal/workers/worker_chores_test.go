package workers

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/prixfixeco/api_server/internal/database"
	mockpublishers "github.com/prixfixeco/api_server/internal/messagequeue/publishers/mock"
	"github.com/prixfixeco/api_server/internal/observability/logging"
)

func TestProvideChoresWorker(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		actual := ProvideChoresWorker(
			logging.NewZerologLogger(),
			&database.MockDatabase{},
			&mockpublishers.Publisher{},
		)
		assert.NotNil(t, actual)
	})
}

func TestChoresWorker_HandleMessage(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		actual := ProvideChoresWorker(
			logging.NewZerologLogger(),
			&database.MockDatabase{},
			&mockpublishers.Publisher{},
		)
		assert.NotNil(t, actual)

		ctx := context.Background()
		assert.NoError(t, actual.HandleMessage(ctx, []byte(`{"choreType":"finalize_meal_plans_with_expired_voting_periods"}`)))
	})

	T.Run("invalid input", func(t *testing.T) {
		t.Parallel()

		actual := ProvideChoresWorker(
			logging.NewZerologLogger(),
			&database.MockDatabase{},
			&mockpublishers.Publisher{},
		)
		assert.NotNil(t, actual)

		ctx := context.Background()
		assert.Error(t, actual.HandleMessage(ctx, []byte("} bad JSON lol")))
	})
}
