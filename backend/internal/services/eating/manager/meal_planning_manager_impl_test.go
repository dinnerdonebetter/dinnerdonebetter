package manager

import (
	"context"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/lib/authentication/sessions"
	msgconfig "github.com/dinnerdonebetter/backend/internal/lib/messagequeue/config"
	mockpublishers "github.com/dinnerdonebetter/backend/internal/lib/messagequeue/mock"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/lib/testutils"
	"github.com/dinnerdonebetter/backend/internal/services/eating/database"
	"github.com/dinnerdonebetter/backend/internal/services/eating/events"
	"github.com/dinnerdonebetter/backend/internal/services/eating/types"
	"github.com/dinnerdonebetter/backend/internal/services/eating/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func matchesEventType(eventType string) any {
	return mock.MatchedBy(func(message *types.DataChangeMessage) bool {
		return message.EventType == eventType
	})
}

func buildMealPlanManagerForTest(t *testing.T) *mealPlanningManager {
	t.Helper()

	queueCfg := &msgconfig.QueuesConfig{
		DataChangesTopicName: "data_changes",
	}

	mpp := &mockpublishers.PublisherProvider{}
	mpp.On("ProvidePublisher", queueCfg.DataChangesTopicName).Return(&mockpublishers.Publisher{}, nil)

	m, err := NewMealPlanningManager(
		logging.NewNoopLogger(),
		tracing.NewNoopTracerProvider(),
		database.NewMockDatabase(),
		queueCfg,
		mpp,
	)
	require.NoError(t, err)

	mock.AssertExpectationsForObjects(t, mpp)

	return m.(*mealPlanningManager)
}

func setupExpectations(
	manager *mealPlanningManager,
	dbSetupFunc func(db *database.MockDatabase),
	eventTypes ...string,
) []any {
	db := database.NewMockDatabase()
	if dbSetupFunc != nil {
		dbSetupFunc(db)
	}
	manager.db = db

	mp := &mockpublishers.Publisher{}
	for _, eventType := range eventTypes {
		mp.On("PublishAsync", testutils.ContextMatcher, matchesEventType(eventType)).Return()
	}
	manager.dataChangesPublisher = mp

	return []any{db, mp}
}

func TestMealPlanningManager_buildDataChangeMessageFromContext(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		mpm := buildMealPlanManagerForTest(t)

		sessionContextData := &sessions.ContextData{
			Requester:         sessions.RequesterInfo{UserID: fakes.BuildFakeID()},
			ActiveHouseholdID: fakes.BuildFakeID(),
		}
		ctx = context.WithValue(ctx, sessions.SessionContextDataKey, sessionContextData)

		expected := &types.DataChangeMessage{
			EventType: events.MealCreated,
			Context: map[string]any{
				"things": "stuff",
			},
			UserID:      sessionContextData.Requester.UserID,
			HouseholdID: sessionContextData.ActiveHouseholdID,
		}

		actual := mpm.buildDataChangeMessageFromContext(ctx, expected.EventType, expected.Context)

		assert.Equal(t, expected, actual)
	})
}

func TestMealPlanningManager_ListMeals(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		mpm := buildMealPlanManagerForTest(t)

		expected := fakes.BuildFakeMealsList()

		expectations := setupExpectations(
			mpm,
			func(db *database.MockDatabase) {
				db.MealDataManagerMock.On("GetMeals", testutils.ContextMatcher, testutils.QueryFilterMatcher).Return(expected, nil)
			},
		)

		actual, cursor, err := mpm.ListMeals(ctx, nil)
		assert.NoError(t, err)
		assert.Empty(t, cursor)
		assert.Equal(t, expected.Data, actual)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestMealPlanningManager_CreateMeal(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		mpm := buildMealPlanManagerForTest(t)

		expected := fakes.BuildFakeMeal()
		fakeInput := fakes.BuildFakeMealCreationRequestInput()

		expectations := setupExpectations(
			mpm,
			func(db *database.MockDatabase) {
				db.MealDataManagerMock.On("CreateMeal", testutils.ContextMatcher, testutils.MatchType[*types.MealDatabaseCreationInput]()).Return(expected, nil)
			},
			events.MealCreated,
		)

		actual, err := mpm.CreateMeal(ctx, fakeInput)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestMealPlanningManager_ReadMeal(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		mpm := buildMealPlanManagerForTest(t)

		expected := fakes.BuildFakeMeal()

		expectations := setupExpectations(
			mpm,
			func(db *database.MockDatabase) {
				db.MealDataManagerMock.On("GetMeal", testutils.ContextMatcher, expected.ID).Return(expected, nil)
			},
		)

		actual, err := mpm.ReadMeal(ctx, expected.ID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestMealPlanningManager_SearchMeals(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		mpm := buildMealPlanManagerForTest(t)

		expected := fakes.BuildFakeMealsList()
		exampleQuery := fakes.BuildFakeID()

		expectations := setupExpectations(
			mpm,
			func(db *database.MockDatabase) {
				db.MealDataManagerMock.On("SearchForMeals", testutils.ContextMatcher, exampleQuery, testutils.QueryFilterMatcher).Return(expected, nil)
			},
		)

		actual, err := mpm.SearchMeals(ctx, exampleQuery, nil)
		assert.NoError(t, err)
		assert.Equal(t, expected.Data, actual)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestMealPlanningManager_ArchiveMeal(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		mpm := buildMealPlanManagerForTest(t)

		expected := fakes.BuildFakeMeal()

		expectations := setupExpectations(
			mpm,
			func(db *database.MockDatabase) {
				db.MealDataManagerMock.On("ArchiveMeal", testutils.ContextMatcher, expected.ID, expected.CreatedByUser).Return(nil)
			},
			events.MealArchived,
		)

		err := mpm.ArchiveMeal(ctx, expected.ID, expected.CreatedByUser)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestMealPlanningManager_ListMealPlans(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		t.SkipNow()
	})
}

func TestMealPlanningManager_CreateMealPlan(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		t.SkipNow()
	})
}

func TestMealPlanningManager_ReadMealPlan(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		t.SkipNow()
	})
}

func TestMealPlanningManager_UpdateMealPlan(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		t.SkipNow()
	})
}

func TestMealPlanningManager_ArchiveMealPlan(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		t.SkipNow()
	})
}

func TestMealPlanningManager_FinalizeMealPlan(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		t.SkipNow()
	})
}

func TestMealPlanningManager_ListMealPlanEvents(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		t.SkipNow()
	})
}

func TestMealPlanningManager_CreateMealPlanEvent(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		t.SkipNow()
	})
}

func TestMealPlanningManager_ReadMealPlanEvent(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		t.SkipNow()
	})
}

func TestMealPlanningManager_UpdateMealPlanEvent(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		t.SkipNow()
	})
}

func TestMealPlanningManager_ArchiveMealPlanEvent(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		t.SkipNow()
	})
}

func TestMealPlanningManager_ListMealPlanOptions(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		t.SkipNow()
	})
}

func TestMealPlanningManager_CreateMealPlanOption(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		t.SkipNow()
	})
}

func TestMealPlanningManager_ReadMealPlanOption(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		t.SkipNow()
	})
}

func TestMealPlanningManager_UpdateMealPlanOption(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		t.SkipNow()
	})
}

func TestMealPlanningManager_ArchiveMealPlanOption(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		t.SkipNow()
	})
}

func TestMealPlanningManager_ListMealPlanOptionVotes(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		t.SkipNow()
	})
}

func TestMealPlanningManager_CreateMealPlanOptionVotes(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		t.SkipNow()
	})
}

func TestMealPlanningManager_ReadMealPlanOptionVote(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		t.SkipNow()
	})
}

func TestMealPlanningManager_UpdateMealPlanOptionVote(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		t.SkipNow()
	})
}

func TestMealPlanningManager_ArchiveMealPlanOptionVote(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		t.SkipNow()
	})
}

func TestMealPlanningManager_ListMealPlanTasksByMealPlan(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		t.SkipNow()
	})
}

func TestMealPlanningManager_ReadMealPlanTask(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		t.SkipNow()
	})
}

func TestMealPlanningManager_CreateMealPlanTask(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		t.SkipNow()
	})
}

func TestMealPlanningManager_MealPlanTaskStatusChange(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		t.SkipNow()
	})
}

func TestMealPlanningManager_ListMealPlanGroceryListItemsByMealPlan(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		t.SkipNow()
	})
}

func TestMealPlanningManager_CreateMealPlanGroceryListItem(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		t.SkipNow()
	})
}

func TestMealPlanningManager_ReadMealPlanGroceryListItem(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		t.SkipNow()
	})
}

func TestMealPlanningManager_UpdateMealPlanGroceryListItem(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		t.SkipNow()
	})
}

func TestMealPlanningManager_ArchiveMealPlanGroceryListItem(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		t.SkipNow()
	})
}

func TestMealPlanningManager_ListIngredientPreferences(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		t.SkipNow()
	})
}

func TestMealPlanningManager_CreateIngredientPreference(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		t.SkipNow()
	})
}

func TestMealPlanningManager_UpdateIngredientPreference(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		t.SkipNow()
	})
}

func TestMealPlanningManager_ArchiveIngredientPreference(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		t.SkipNow()
	})
}

func TestMealPlanningManager_ListInstrumentOwnerships(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		t.SkipNow()
	})
}

func TestMealPlanningManager_CreateInstrumentOwnership(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		t.SkipNow()
	})
}

func TestMealPlanningManager_ReadInstrumentOwnership(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		t.SkipNow()
	})
}

func TestMealPlanningManager_UpdateInstrumentOwnership(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		t.SkipNow()
	})
}

func TestMealPlanningManager_ArchiveInstrumentOwnership(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		t.SkipNow()
	})
}
