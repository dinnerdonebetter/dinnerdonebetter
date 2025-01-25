package workers

import (
	"context"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/database"
	mockpublishers "github.com/dinnerdonebetter/backend/internal/lib/messagequeue/mock"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/metrics"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/testutils"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func newTestChoresWorker(t *testing.T) *mealPlanFinalizationWorker {
	t.Helper()

	worker, err := ProvideMealPlanFinalizationWorker(
		logging.NewNoopLogger(),
		&database.MockDatabase{},
		&mockpublishers.Publisher{},
		tracing.NewNoopTracerProvider(),
		metrics.NewNoopMetricsProvider(),
	)
	assert.NotNil(t, worker)
	assert.NoError(t, err)

	return worker.(*mealPlanFinalizationWorker)
}

func TestProvideChoresWorker(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		actual, err := ProvideMealPlanFinalizationWorker(
			logging.NewNoopLogger(),
			&database.MockDatabase{},
			&mockpublishers.Publisher{},
			tracing.NewNoopTracerProvider(),
			metrics.NewNoopMetricsProvider(),
		)
		assert.NotNil(t, actual)
		assert.NoError(t, err)
	})
}

func TestChoresWorker_FinalizeExpiredMealPlansWithoutReturningCount(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleMealPlans := fakes.BuildFakeMealPlansList().Data

		dbm := database.NewMockDatabase()
		dbm.MealPlanDataManagerMock.On(
			"GetUnfinalizedMealPlansWithExpiredVotingPeriods",
			testutils.ContextMatcher,
		).Return(exampleMealPlans, nil)

		pup := &mockpublishers.Publisher{}

		for _, mealPlan := range exampleMealPlans {
			dbm.MealPlanDataManagerMock.On(
				"AttemptToFinalizeMealPlan",
				testutils.ContextMatcher,
				mealPlan.ID,
				mealPlan.BelongsToHousehold,
			).Return(true, nil)

			pup.On("Publish", testutils.ContextMatcher, mock.AnythingOfType("*types.DataChangeMessage")).Return(nil)
		}

		worker := newTestChoresWorker(t)
		worker.dataManager = dbm
		worker.postUpdatesPublisher = pup

		assert.NoError(t, worker.FinalizeExpiredMealPlansWithoutReturningCount(ctx))

		mock.AssertExpectationsForObjects(t, dbm, pup)
	})
}
