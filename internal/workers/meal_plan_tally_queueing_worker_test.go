package workers

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/prixfixeco/backend/internal/analytics"
	"github.com/prixfixeco/backend/internal/database"
	"github.com/prixfixeco/backend/internal/email"
	mockpublishers "github.com/prixfixeco/backend/internal/messagequeue/mock"
	"github.com/prixfixeco/backend/internal/observability/logging"
	"github.com/prixfixeco/backend/internal/observability/logging/zerolog"
	"github.com/prixfixeco/backend/internal/observability/tracing"
	"github.com/prixfixeco/backend/pkg/types/fakes"
	testutils "github.com/prixfixeco/backend/tests/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func newTestTallyQueueingWorker(t *testing.T) *MealPlanTallyScheduler {
	t.Helper()

	worker := ProvideMealPlanTallyScheduler(
		zerolog.NewZerologLogger(logging.DebugLevel),
		&database.MockDatabase{},
		&mockpublishers.Publisher{},
		&email.MockEmailer{},
		&analytics.MockEventReporter{},
		tracing.NewNoopTracerProvider(),
	)
	assert.NotNil(t, worker)

	return worker
}

func TestProvideMealPlanTallyQueueingWorker(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		actual := ProvideMealPlanTallyScheduler(
			zerolog.NewZerologLogger(logging.DebugLevel),
			&database.MockDatabase{},
			&mockpublishers.Publisher{},
			&email.MockEmailer{},
			&analytics.MockEventReporter{},
			tracing.NewNoopTracerProvider(),
		)
		assert.NotNil(t, actual)
	})
}

func TestMealPlanTallyQueueingWorker_FinalizeExpiredMealPlansWithoutReturningCount(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		body, err := json.Marshal(&struct{}{})
		require.NoError(t, err)

		exampleMealPlans := fakes.BuildFakeMealPlanList().Data

		dbm := database.NewMockDatabase()
		dbm.MealPlanDataManager.On(
			"GetUnfinalizedMealPlansWithExpiredVotingPeriods",
			testutils.ContextMatcher,
		).Return(exampleMealPlans, nil)

		mqm := &mockpublishers.Publisher{}

		for _, mealPlan := range exampleMealPlans {
			mqm.On("Publish", testutils.ContextMatcher, &MealPlanTallyRequest{
				MealPlanID:  mealPlan.ID,
				HouseholdID: mealPlan.BelongsToHousehold,
			}).Return(nil)
		}

		worker := newTestTallyQueueingWorker(t)
		worker.dataManager = dbm
		worker.tallyQueuePublisher = mqm

		assert.NoError(t, worker.ScheduleMealPlanTallies(ctx, body))

		mock.AssertExpectationsForObjects(t, dbm, mqm)
	})
}
