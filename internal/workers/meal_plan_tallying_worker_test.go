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

func newTestMealPlanTallyingWorker(t *testing.T) *MealPlanTallyingWorker {
	t.Helper()

	worker := ProvideMealPlanTallyingWorker(
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

func TestMealPlanTallyingWorker_TallyMealPlanVotes(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleMealPlan := fakes.BuildFakeMealPlan()

		exampleInput := &MealPlanTallyRequest{
			MealPlanID:  exampleMealPlan.ID,
			HouseholdID: exampleMealPlan.BelongsToHousehold,
		}
		body, err := json.Marshal(exampleInput)
		require.NoError(t, err)

		dbm := database.NewMockDatabase()
		dbm.MealPlanDataManager.On(
			"AttemptToFinalizeMealPlan",
			testutils.ContextMatcher,
			exampleInput.MealPlanID,
			exampleInput.HouseholdID,
		).Return(true, nil)

		worker := newTestMealPlanTallyingWorker(t)
		worker.dataManager = dbm

		assert.NoError(t, worker.TallyMealPlanVotes(ctx, body))

		mock.AssertExpectationsForObjects(t, dbm)
	})
}
