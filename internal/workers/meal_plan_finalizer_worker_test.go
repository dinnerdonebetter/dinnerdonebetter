package workers

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/prixfixeco/api_server/internal/customerdata"
	"github.com/prixfixeco/api_server/internal/database"
	"github.com/prixfixeco/api_server/internal/email"
	mockpublishers "github.com/prixfixeco/api_server/internal/messagequeue/mock"
	"github.com/prixfixeco/api_server/internal/observability/logging/zerolog"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
	"github.com/prixfixeco/api_server/pkg/types/fakes"
	testutils "github.com/prixfixeco/api_server/tests/utils"
)

func newTestChoresWorker(t *testing.T) *MealPlanFinalizationWorker {
	t.Helper()

	worker := ProvideMealPlanFinalizationWorker(
		zerolog.NewZerologLogger(),
		&database.MockDatabase{},
		&mockpublishers.Publisher{},
		&email.MockEmailer{},
		&customerdata.MockCollector{},
		tracing.NewNoopTracerProvider(),
	)
	assert.NotNil(t, worker)

	return worker
}

func TestProvideChoresWorker(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		actual := ProvideMealPlanFinalizationWorker(
			zerolog.NewZerologLogger(),
			&database.MockDatabase{},
			&mockpublishers.Publisher{},
			&email.MockEmailer{},
			&customerdata.MockCollector{},
			tracing.NewNoopTracerProvider(),
		)
		assert.NotNil(t, actual)
	})
}

func TestChoresWorker_FinalizeExpiredMealPlansWithoutReturningCount(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleInput := &types.ChoreMessage{
			ChoreType:                 types.FinalizeMealPlansWithExpiredVotingPeriodsChoreType,
			MealPlanID:                fakes.BuildFakeID(),
			AttributableToHouseholdID: fakes.BuildFakeID(),
		}
		body, err := json.Marshal(exampleInput)
		require.NoError(t, err)

		exampleMealPlans := fakes.BuildFakeMealPlanList().MealPlans

		dbm := database.NewMockDatabase()
		dbm.MealPlanDataManager.On(
			"GetUnfinalizedMealPlansWithExpiredVotingPeriods",
			testutils.ContextMatcher,
		).Return(exampleMealPlans, nil)

		mqm := &mockpublishers.Publisher{}

		for _, mealPlan := range exampleMealPlans {
			dbm.MealPlanDataManager.On(
				"AttemptToFinalizeMealPlan",
				testutils.ContextMatcher,
				mealPlan.ID,
				mealPlan.BelongsToHousehold,
			).Return(true, nil)

			mqm.On(
				"Publish",
				testutils.ContextMatcher,
				mock.AnythingOfType("*types.DataChangeMessage"),
			).Return(nil)
		}

		worker := newTestChoresWorker(t)
		worker.dataManager = dbm
		worker.postUpdatesPublisher = mqm

		assert.NoError(t, worker.FinalizeExpiredMealPlansWithoutReturningCount(ctx, body))

		mock.AssertExpectationsForObjects(t, dbm, mqm)
	})
}
