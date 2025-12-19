package mealplanfinalizer

import (
	"testing"

	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning/fakes"
	mealplanningmock "github.com/dinnerdonebetter/backend/internal/domain/mealplanning/mocks"
	msgconfig "github.com/dinnerdonebetter/backend/internal/platform/messagequeue/config"
	mockpublishers "github.com/dinnerdonebetter/backend/internal/platform/messagequeue/mock"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/metrics"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/platform/testutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func buildNewMealPlanFinalizerForTest(t *testing.T) *Worker {
	t.Helper()

	ctx := t.Context()
	cfg := &msgconfig.QueuesConfig{DataChangesTopicName: "data_changes"}

	pp := &mockpublishers.PublisherProvider{}
	pp.On("ProvidePublisher", cfg.DataChangesTopicName).Return(&mockpublishers.Publisher{}, nil)

	x, err := NewMealPlanFinalizer(
		ctx,
		logging.NewNoopLogger(),
		tracing.NewNoopTracerProvider(),
		&mealplanningmock.Repository{},
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

		dbm := &mealplanningmock.Repository{}
		dbm.On(
			"GetUnfinalizedMealPlansWithExpiredVotingPeriods",
			testutils.ContextMatcher,
		).Return(exampleMealPlans, nil)

		pup := &mockpublishers.Publisher{}

		for _, mealPlan := range exampleMealPlans {
			dbm.On(
				"AttemptToFinalizeMealPlan",
				testutils.ContextMatcher,
				mealPlan.ID,
				mealPlan.BelongsToAccount,
			).Return(true, nil)

			pup.On("Publish", testutils.ContextMatcher, mock.AnythingOfType("*audit.DataChangeMessage")).Return(nil)
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
