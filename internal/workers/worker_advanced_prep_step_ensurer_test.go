package workers

import (
	"context"
	testutils "github.com/prixfixeco/api_server/tests/utils"
	"github.com/stretchr/testify/mock"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/prixfixeco/api_server/internal/customerdata"
	"github.com/prixfixeco/api_server/internal/database"
	"github.com/prixfixeco/api_server/internal/graphing"
	mockpublishers "github.com/prixfixeco/api_server/internal/messagequeue/mock"
	"github.com/prixfixeco/api_server/internal/observability/logging/zerolog"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
)

func TestProvideAdvancedPrepStepCreationEnsurerWorker(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		actual := ProvideAdvancedPrepStepCreationEnsurerWorker(
			zerolog.NewZerologLogger(),
			&database.MockDatabase{},
			&graphing.MockRecipeGrapher{},
			&mockpublishers.Publisher{},
			&customerdata.MockCollector{},
			tracing.NewNoopTracerProvider(),
		)
		assert.NotNil(t, actual)
	})
}

func TestAdvancedPrepStepCreationEnsurerWorker_HandleMessage(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		w := ProvideAdvancedPrepStepCreationEnsurerWorker(
			zerolog.NewZerologLogger(),
			&database.MockDatabase{},
			&graphing.MockRecipeGrapher{},
			&mockpublishers.Publisher{},
			&customerdata.MockCollector{},
			tracing.NewNoopTracerProvider(),
		)
		assert.NotNil(t, w)

		ctx := context.Background()

		mdm := database.NewMockDatabase()
		mdm.MealPlanDataManager.On("GetFinalizedMealPlanIDsForTheNextWeek", testutils.ContextMatcher).Return([]*types.FinalizedMealPlanDatabaseResult{}, nil)
		w.dataManager = mdm

		err := w.HandleMessage(ctx, []byte("{}"))
		assert.NoError(t, err)
	})
}

func TestAdvancedPrepStepCreationEnsurerWorker_DetermineCreatableSteps(T *testing.T) {
	T.Parallel()

	T.Run("with nothing to do", func(t *testing.T) {
		t.Parallel()

		w := ProvideAdvancedPrepStepCreationEnsurerWorker(
			zerolog.NewZerologLogger(),
			&database.MockDatabase{},
			&graphing.MockRecipeGrapher{},
			&mockpublishers.Publisher{},
			&customerdata.MockCollector{},
			tracing.NewNoopTracerProvider(),
		)
		assert.NotNil(t, w)

		ctx := context.Background()
		expected := []*types.AdvancedPrepStepDatabaseCreationInput{}

		mdm := database.NewMockDatabase()
		mdm.MealPlanDataManager.On("GetFinalizedMealPlanIDsForTheNextWeek", testutils.ContextMatcher).Return([]*types.FinalizedMealPlanDatabaseResult{}, nil)
		w.dataManager = mdm

		actual, err := w.DetermineCreatableSteps(ctx)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mock.AssertExpectationsForObjects(t, mdm)
	})

	T.Run("standard", func(t *testing.T) {
		t.SkipNow()

		w := ProvideAdvancedPrepStepCreationEnsurerWorker(
			zerolog.NewZerologLogger(),
			&database.MockDatabase{},
			&graphing.MockRecipeGrapher{},
			&mockpublishers.Publisher{},
			&customerdata.MockCollector{},
			tracing.NewNoopTracerProvider(),
		)
		assert.NotNil(t, w)

		ctx := context.Background()
		expected := []*types.AdvancedPrepStepDatabaseCreationInput{}

		actual, err := w.DetermineCreatableSteps(ctx)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})
}
