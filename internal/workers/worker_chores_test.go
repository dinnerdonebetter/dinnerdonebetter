package workers

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace"

	"github.com/prixfixeco/api_server/internal/customerdata"
	"github.com/prixfixeco/api_server/internal/database"
	"github.com/prixfixeco/api_server/internal/email"
	mockpublishers "github.com/prixfixeco/api_server/internal/messagequeue/mock"
	"github.com/prixfixeco/api_server/internal/observability/logging/zerolog"
	"github.com/prixfixeco/api_server/pkg/types"
	"github.com/prixfixeco/api_server/pkg/types/fakes"
	testutils "github.com/prixfixeco/api_server/tests/utils"
)

func TestProvideChoresWorker(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		actual := ProvideChoresWorker(
			zerolog.NewZerologLogger(),
			&database.MockDatabase{},
			&mockpublishers.Publisher{},
			&email.MockEmailer{},
			&customerdata.MockCollector{},
			trace.NewNoopTracerProvider(),
		)
		assert.NotNil(t, actual)
	})
}

func TestChoresWorker_HandleMessage(T *testing.T) {
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

		dbm := database.NewMockDatabase()
		dbm.MealPlanDataManager.On(
			"FinalizeMealPlanWithExpiredVotingPeriod",
			testutils.ContextMatcher,
			exampleInput.MealPlanID,
			exampleInput.AttributableToHouseholdID,
		).Return(true, nil)

		worker := newTestChoresWorker(t)
		worker.dataManager = dbm

		assert.NoError(t, worker.HandleMessage(ctx, body))

		mock.AssertExpectationsForObjects(t, dbm)
	})

	T.Run("invalid input", func(t *testing.T) {
		t.Parallel()

		worker := newTestChoresWorker(t)

		ctx := context.Background()
		assert.Error(t, worker.HandleMessage(ctx, []byte("} bad JSON lol")))
	})
}
