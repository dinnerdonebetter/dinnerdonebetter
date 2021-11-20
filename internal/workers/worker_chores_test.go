package workers

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/prixfixeco/api_server/internal/database"
	mockpublishers "github.com/prixfixeco/api_server/internal/messagequeue/publishers/mock"
	"github.com/prixfixeco/api_server/internal/observability/logging"
	"github.com/prixfixeco/api_server/pkg/types"
	"github.com/prixfixeco/api_server/pkg/types/fakes"
	testutils "github.com/prixfixeco/api_server/tests/utils"
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

		dbm := database.NewMockDatabase()

		actual := ProvideChoresWorker(
			logging.NewZerologLogger(),
			dbm,
			&mockpublishers.Publisher{},
		)
		assert.NotNil(t, actual)

		ctx := context.Background()
		exampleInput := &types.ChoreMessage{
			ChoreType:                 types.FinalizeMealPlansWithExpiredVotingPeriodsChoreType,
			MealPlanID:                fakes.BuildFakeID(),
			AttributableToHouseholdID: fakes.BuildFakeID(),
		}
		body, err := json.Marshal(exampleInput)
		require.NoError(t, err)

		dbm.MealPlanDataManager.On(
			"FinalizeMealPlan",
			testutils.ContextMatcher,
			exampleInput.MealPlanID,
			exampleInput.AttributableToHouseholdID,
			false,
		).Return(true, nil)

		assert.NoError(t, actual.HandleMessage(ctx, body))

		mock.AssertExpectationsForObjects(t, dbm)
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
