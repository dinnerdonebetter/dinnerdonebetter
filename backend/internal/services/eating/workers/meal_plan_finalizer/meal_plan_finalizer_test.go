package mealplanfinalizer

import (
	"testing"

	"github.com/dinnerdonebetter/backend/internal/database"
	msgconfig "github.com/dinnerdonebetter/backend/internal/lib/messagequeue/config"
	mockpublishers "github.com/dinnerdonebetter/backend/internal/lib/messagequeue/mock"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/metrics"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/lib/testutils"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func buildNewMealPlanFinalizerForTest(t *testing.T) *Worker {
	t.Helper()

	cfg := &msgconfig.QueuesConfig{DataChangesTopicName: "data_changes"}

	pp := &mockpublishers.PublisherProvider{}
	pp.On("ProvidePublisher", cfg.DataChangesTopicName).Return(&mockpublishers.Publisher{}, nil)

	x, err := NewMealPlanFinalizer(
		logging.NewNoopLogger(),
		tracing.NewNoopTracerProvider(),
		database.NewMockDatabase(),
		pp,
		metrics.NewNoopMetricsProvider(),
		cfg,
	)
	require.NoError(t, err)

	return x
}

func TestWorker_Work(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
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

		worker := buildNewMealPlanFinalizerForTest(t)
		worker.dataManager = dbm
		worker.postUpdatesPublisher = pup

		expected := int64(len(exampleMealPlans))

		actual, err := worker.Work(ctx)
		assert.Equal(t, expected, actual)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, dbm, pup)
	})
}
