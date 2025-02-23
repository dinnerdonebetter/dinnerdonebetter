package manager

import (
	"context"
	"slices"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/lib/authentication/sessions"
	msgconfig "github.com/dinnerdonebetter/backend/internal/lib/messagequeue/config"
	mockpublishers "github.com/dinnerdonebetter/backend/internal/lib/messagequeue/mock"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/keys"
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

func eventMatches(eventType string, keys []string) any {
	return mock.MatchedBy(func(message *types.DataChangeMessage) bool {
		allContextKeys := []string{}
		for k := range message.Context {
			allContextKeys = append(allContextKeys, k)
		}

		return slices.Equal(keys, allContextKeys) && message.EventType == eventType
	})
}

func setupExpectations(
	manager *mealPlanningManager,
	dbSetupFunc func(db *database.MockDatabase),
	eventTypeMaps ...map[string][]string,
) []any {
	db := database.NewMockDatabase()
	if dbSetupFunc != nil {
		dbSetupFunc(db)
	}
	manager.db = db

	mp := &mockpublishers.Publisher{}
	for _, eventTypeMap := range eventTypeMaps {
		for eventType, payload := range eventTypeMap {
			mp.On("PublishAsync", testutils.ContextMatcher, eventMatches(eventType, payload)).Return()
		}
	}
	manager.dataChangesPublisher = mp

	return []any{db, mp}
}

func buildMealPlanManagerForTest(t *testing.T) *mealPlanningManager {
	t.Helper()

	queueCfg := &msgconfig.QueuesConfig{
		DataChangesTopicName: t.Name(),
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
				db.MealDataManagerMock.On(testutils.GetMethodName(mpm.db.GetMeals), testutils.ContextMatcher, testutils.QueryFilterMatcher).Return(expected, nil)
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
				db.MealDataManagerMock.On(testutils.GetMethodName(mpm.db.CreateMeal), testutils.ContextMatcher, testutils.MatchType[*types.MealDatabaseCreationInput]()).Return(expected, nil)
			},
			map[string][]string{
				events.MealCreated: {keys.MealIDKey},
			},
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
				db.MealDataManagerMock.On(testutils.GetMethodName(mpm.db.GetMeal), testutils.ContextMatcher, expected.ID).Return(expected, nil)
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
				db.MealDataManagerMock.On(testutils.GetMethodName(mpm.db.SearchForMeals), testutils.ContextMatcher, exampleQuery, testutils.QueryFilterMatcher).Return(expected, nil)
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
				db.MealDataManagerMock.On(testutils.GetMethodName(mpm.db.ArchiveMeal), testutils.ContextMatcher, expected.ID, expected.CreatedByUser).Return(nil)
			},
			map[string][]string{
				events.MealArchived: {keys.MealIDKey},
			},
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

		ctx := t.Context()
		mpm := buildMealPlanManagerForTest(t)

		expected := fakes.BuildFakeMealPlansList()
		exampleOwnerID := fakes.BuildFakeID()

		expectations := setupExpectations(
			mpm,
			func(db *database.MockDatabase) {
				db.MealPlanDataManagerMock.On(testutils.GetMethodName(mpm.db.GetMealPlansForHousehold), testutils.ContextMatcher, exampleOwnerID, testutils.QueryFilterMatcher).Return(expected, nil)
			},
		)

		actual, cursor, err := mpm.ListMealPlans(ctx, exampleOwnerID, nil)
		assert.NoError(t, err)
		assert.Empty(t, cursor)
		assert.Equal(t, expected.Data, actual)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestMealPlanningManager_CreateMealPlan(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		mpm := buildMealPlanManagerForTest(t)

		expected := fakes.BuildFakeMealPlan()
		fakeInput := fakes.BuildFakeMealPlanCreationRequestInput()

		expectations := setupExpectations(
			mpm,
			func(db *database.MockDatabase) {
				db.MealPlanDataManagerMock.On(testutils.GetMethodName(mpm.db.CreateMealPlan), testutils.ContextMatcher, testutils.MatchType[*types.MealPlanDatabaseCreationInput]()).Return(expected, nil)
			},
			map[string][]string{
				events.MealPlanCreated: {keys.MealPlanIDKey},
			},
		)

		actual, err := mpm.CreateMealPlan(ctx, fakeInput)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestMealPlanningManager_ReadMealPlan(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		mpm := buildMealPlanManagerForTest(t)

		exampleMealPlanID := fakes.BuildFakeID()
		expected := fakes.BuildFakeMealPlan()

		expectations := setupExpectations(
			mpm,
			func(db *database.MockDatabase) {
				db.MealPlanDataManagerMock.On(testutils.GetMethodName(mpm.db.GetMealPlan), testutils.ContextMatcher, exampleMealPlanID, expected.ID).Return(expected, nil)
			},
		)

		actual, err := mpm.ReadMealPlan(ctx, exampleMealPlanID, expected.ID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mock.AssertExpectationsForObjects(t, expectations...)
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

		ctx := t.Context()
		mpm := buildMealPlanManagerForTest(t)

		expected := fakes.BuildFakeMealPlanEventsList()
		exampleMealPlanID := fakes.BuildFakeID()

		expectations := setupExpectations(
			mpm,
			func(db *database.MockDatabase) {
				db.MealPlanEventDataManagerMock.On(testutils.GetMethodName(mpm.db.GetMealPlanEvents), testutils.ContextMatcher, exampleMealPlanID, testutils.QueryFilterMatcher).Return(expected, nil)
			},
		)

		actual, cursor, err := mpm.ListMealPlanEvents(ctx, exampleMealPlanID, nil)
		assert.NoError(t, err)
		assert.Empty(t, cursor)
		assert.Equal(t, expected.Data, actual)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestMealPlanningManager_CreateMealPlanEvent(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		mpm := buildMealPlanManagerForTest(t)

		expected := fakes.BuildFakeMealPlanEvent()
		fakeInput := fakes.BuildFakeMealPlanEventCreationRequestInput()

		expectations := setupExpectations(
			mpm,
			func(db *database.MockDatabase) {
				db.MealPlanEventDataManagerMock.On(testutils.GetMethodName(mpm.db.CreateMealPlanEvent), testutils.ContextMatcher, testutils.MatchType[*types.MealPlanEventDatabaseCreationInput]()).Return(expected, nil)
			},
			map[string][]string{
				events.MealPlanEventCreated: {keys.MealPlanEventIDKey},
			},
		)

		actual, err := mpm.CreateMealPlanEvent(ctx, fakeInput)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestMealPlanningManager_ReadMealPlanEvent(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		mpm := buildMealPlanManagerForTest(t)

		exampleMealPlanID := fakes.BuildFakeID()
		expected := fakes.BuildFakeMealPlanEvent()

		expectations := setupExpectations(
			mpm,
			func(db *database.MockDatabase) {
				db.MealPlanEventDataManagerMock.On(testutils.GetMethodName(mpm.db.GetMealPlanEvent), testutils.ContextMatcher, exampleMealPlanID, expected.ID).Return(expected, nil)
			},
		)

		actual, err := mpm.ReadMealPlanEvent(ctx, exampleMealPlanID, expected.ID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mock.AssertExpectationsForObjects(t, expectations...)
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

		ctx := t.Context()
		mpm := buildMealPlanManagerForTest(t)

		expected := fakes.BuildFakeMealPlanOptionsList()
		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanEventID := fakes.BuildFakeID()

		expectations := setupExpectations(
			mpm,
			func(db *database.MockDatabase) {
				db.MealPlanOptionDataManagerMock.On(testutils.GetMethodName(mpm.db.GetMealPlanOptions), testutils.ContextMatcher, exampleMealPlanID, exampleMealPlanEventID, testutils.QueryFilterMatcher).Return(expected, nil)
			},
		)

		actual, cursor, err := mpm.ListMealPlanOptions(ctx, exampleMealPlanID, exampleMealPlanEventID, nil)
		assert.NoError(t, err)
		assert.Empty(t, cursor)
		assert.Equal(t, expected.Data, actual)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestMealPlanningManager_CreateMealPlanOption(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		mpm := buildMealPlanManagerForTest(t)

		expected := fakes.BuildFakeMealPlanOption()
		fakeInput := fakes.BuildFakeMealPlanOptionCreationRequestInput()

		expectations := setupExpectations(
			mpm,
			func(db *database.MockDatabase) {
				db.MealPlanOptionDataManagerMock.On(testutils.GetMethodName(mpm.db.CreateMealPlanOption), testutils.ContextMatcher, testutils.MatchType[*types.MealPlanOptionDatabaseCreationInput]()).Return(expected, nil)
			},
			map[string][]string{
				events.MealPlanOptionCreated: {keys.MealPlanOptionIDKey},
			},
		)

		actual, err := mpm.CreateMealPlanOption(ctx, fakeInput)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestMealPlanningManager_ReadMealPlanOption(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		mpm := buildMealPlanManagerForTest(t)

		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanEventID := fakes.BuildFakeID()
		expected := fakes.BuildFakeMealPlanOption()

		expectations := setupExpectations(
			mpm,
			func(db *database.MockDatabase) {
				db.MealPlanOptionDataManagerMock.On(testutils.GetMethodName(mpm.db.GetMealPlanOption), testutils.ContextMatcher, exampleMealPlanID, exampleMealPlanEventID, expected.ID).Return(expected, nil)
			},
		)

		actual, err := mpm.ReadMealPlanOption(ctx, exampleMealPlanID, exampleMealPlanEventID, expected.ID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mock.AssertExpectationsForObjects(t, expectations...)
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

		ctx := t.Context()
		mpm := buildMealPlanManagerForTest(t)

		expected := fakes.BuildFakeMealPlanOptionVotesList()
		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanEventID := fakes.BuildFakeID()
		exampleMealPlanOptionID := fakes.BuildFakeID()

		expectations := setupExpectations(
			mpm,
			func(db *database.MockDatabase) {
				db.MealPlanOptionVoteDataManagerMock.On(testutils.GetMethodName(mpm.db.GetMealPlanOptionVotes), testutils.ContextMatcher, exampleMealPlanID, exampleMealPlanEventID, exampleMealPlanOptionID, testutils.QueryFilterMatcher).Return(expected, nil)
			},
		)

		actual, cursor, err := mpm.ListMealPlanOptionVotes(ctx, exampleMealPlanID, exampleMealPlanEventID, exampleMealPlanOptionID, nil)
		assert.NoError(t, err)
		assert.Empty(t, cursor)
		assert.Equal(t, expected.Data, actual)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestMealPlanningManager_CreateMealPlanOptionVotes(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		mpm := buildMealPlanManagerForTest(t)

		expected := fakes.BuildFakeMealPlanOptionVotesList().Data
		fakeInput := fakes.BuildFakeMealPlanOptionVoteCreationRequestInput()

		expectations := setupExpectations(
			mpm,
			func(db *database.MockDatabase) {
				db.MealPlanOptionVoteDataManagerMock.On(testutils.GetMethodName(mpm.db.CreateMealPlanOptionVote), testutils.ContextMatcher, testutils.MatchType[*types.MealPlanOptionVotesDatabaseCreationInput]()).Return(expected, nil)
			},
			map[string][]string{
				events.MealPlanOptionVoteCreated: {"vote_count", "created"},
			},
		)

		actual, err := mpm.CreateMealPlanOptionVotes(ctx, fakeInput)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestMealPlanningManager_ReadMealPlanOptionVote(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		mpm := buildMealPlanManagerForTest(t)

		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanEventID := fakes.BuildFakeID()
		exampleMealPlanOptionID := fakes.BuildFakeID()
		expected := fakes.BuildFakeMealPlanOptionVote()

		expectations := setupExpectations(
			mpm,
			func(db *database.MockDatabase) {
				db.MealPlanOptionVoteDataManagerMock.On(testutils.GetMethodName(mpm.db.GetMealPlanOptionVote), testutils.ContextMatcher, exampleMealPlanID, exampleMealPlanEventID, exampleMealPlanOptionID, expected.ID).Return(expected, nil)
			},
		)

		actual, err := mpm.ReadMealPlanOptionVote(ctx, exampleMealPlanID, exampleMealPlanEventID, exampleMealPlanOptionID, expected.ID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mock.AssertExpectationsForObjects(t, expectations...)
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

		ctx := t.Context()
		mpm := buildMealPlanManagerForTest(t)

		expected := fakes.BuildFakeMealPlanTasksList()
		exampleMealPlanID := fakes.BuildFakeID()

		expectations := setupExpectations(
			mpm,
			func(db *database.MockDatabase) {
				db.MealPlanTaskDataManagerMock.On(testutils.GetMethodName(mpm.db.GetMealPlanTasksForMealPlan), testutils.ContextMatcher, exampleMealPlanID).Return(expected.Data, nil)
			},
		)

		actual, cursor, err := mpm.ListMealPlanTasksByMealPlan(ctx, exampleMealPlanID, nil)
		assert.NoError(t, err)
		assert.Empty(t, cursor)
		assert.Equal(t, expected.Data, actual)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestMealPlanningManager_ReadMealPlanTask(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		mpm := buildMealPlanManagerForTest(t)

		exampleMealPlanID := fakes.BuildFakeID()
		expected := fakes.BuildFakeMealPlanTask()

		expectations := setupExpectations(
			mpm,
			func(db *database.MockDatabase) {
				db.MealPlanTaskDataManagerMock.On(testutils.GetMethodName(mpm.db.GetMealPlanTask), testutils.ContextMatcher, exampleMealPlanID, expected.ID).Return(expected, nil)
			},
		)

		actual, err := mpm.ReadMealPlanTask(ctx, exampleMealPlanID, expected.ID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestMealPlanningManager_CreateMealPlanTask(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		mpm := buildMealPlanManagerForTest(t)

		expected := fakes.BuildFakeMealPlanTask()
		fakeInput := fakes.BuildFakeMealPlanTaskCreationRequestInput()

		expectations := setupExpectations(
			mpm,
			func(db *database.MockDatabase) {
				db.MealPlanTaskDataManagerMock.On(testutils.GetMethodName(mpm.db.CreateMealPlanTask), testutils.ContextMatcher, testutils.MatchType[*types.MealPlanTaskDatabaseCreationInput]()).Return(expected, nil)
			},
			map[string][]string{
				events.MealPlanTaskCreated: {keys.MealPlanTaskIDKey},
			},
		)

		actual, err := mpm.CreateMealPlanTask(ctx, fakeInput)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mock.AssertExpectationsForObjects(t, expectations...)
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

		ctx := t.Context()
		mpm := buildMealPlanManagerForTest(t)

		expected := fakes.BuildFakeMealPlanGroceryListItemsList()
		exampleMealPlanID := fakes.BuildFakeID()

		expectations := setupExpectations(
			mpm,
			func(db *database.MockDatabase) {
				db.MealPlanGroceryListItemDataManagerMock.On(testutils.GetMethodName(mpm.db.GetMealPlanGroceryListItemsForMealPlan), testutils.ContextMatcher, exampleMealPlanID).Return(expected.Data, nil)
			},
		)

		actual, cursor, err := mpm.ListMealPlanGroceryListItemsByMealPlan(ctx, exampleMealPlanID, nil)
		assert.NoError(t, err)
		assert.Empty(t, cursor)
		assert.Equal(t, expected.Data, actual)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestMealPlanningManager_CreateMealPlanGroceryListItem(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		mpm := buildMealPlanManagerForTest(t)

		expected := fakes.BuildFakeMealPlanGroceryListItem()
		fakeInput := fakes.BuildFakeMealPlanGroceryListItemCreationRequestInput()

		expectations := setupExpectations(
			mpm,
			func(db *database.MockDatabase) {
				db.MealPlanGroceryListItemDataManagerMock.On(testutils.GetMethodName(mpm.db.CreateMealPlanGroceryListItem), testutils.ContextMatcher, testutils.MatchType[*types.MealPlanGroceryListItemDatabaseCreationInput]()).Return(expected, nil)
			},
			map[string][]string{
				events.MealPlanGroceryListItemCreated: {keys.MealPlanGroceryListItemIDKey},
			},
		)

		actual, err := mpm.CreateMealPlanGroceryListItem(ctx, fakeInput)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestMealPlanningManager_ReadMealPlanGroceryListItem(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		mpm := buildMealPlanManagerForTest(t)

		exampleMealPlanID := fakes.BuildFakeID()
		expected := fakes.BuildFakeMealPlanGroceryListItem()

		expectations := setupExpectations(
			mpm,
			func(db *database.MockDatabase) {
				db.MealPlanGroceryListItemDataManagerMock.On(testutils.GetMethodName(mpm.db.GetMealPlanGroceryListItem), testutils.ContextMatcher, exampleMealPlanID, expected.ID).Return(expected, nil)
			},
		)

		actual, err := mpm.ReadMealPlanGroceryListItem(ctx, exampleMealPlanID, expected.ID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mock.AssertExpectationsForObjects(t, expectations...)
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

		ctx := t.Context()
		mpm := buildMealPlanManagerForTest(t)

		expected := fakes.BuildFakeIngredientPreferencesList()
		exampleOwnerID := fakes.BuildFakeID()

		expectations := setupExpectations(
			mpm,
			func(db *database.MockDatabase) {
				db.IngredientPreferenceDataManagerMock.On(testutils.GetMethodName(mpm.db.GetIngredientPreferences), testutils.ContextMatcher, exampleOwnerID, testutils.QueryFilterMatcher).Return(expected, nil)
			},
		)

		actual, cursor, err := mpm.ListIngredientPreferences(ctx, exampleOwnerID, nil)
		assert.NoError(t, err)
		assert.Empty(t, cursor)
		assert.Equal(t, expected.Data, actual)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestMealPlanningManager_CreateIngredientPreference(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		mpm := buildMealPlanManagerForTest(t)

		expected := fakes.BuildFakeIngredientPreferencesList().Data
		fakeInput := fakes.BuildFakeIngredientPreferenceCreationRequestInput()

		expectations := setupExpectations(
			mpm,
			func(db *database.MockDatabase) {
				db.IngredientPreferenceDataManagerMock.On(testutils.GetMethodName(mpm.db.CreateIngredientPreference), testutils.ContextMatcher, testutils.MatchType[*types.IngredientPreferenceDatabaseCreationInput]()).Return(expected, nil)
			},
			map[string][]string{
				events.IngredientPreferenceCreated: {keys.ValidIngredientGroupIDKey, keys.ValidIngredientIDKey, "created"},
			},
		)

		actual, err := mpm.CreateIngredientPreference(ctx, fakeInput)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mock.AssertExpectationsForObjects(t, expectations...)
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

		ctx := t.Context()
		mpm := buildMealPlanManagerForTest(t)

		expected := fakes.BuildFakeInstrumentOwnershipsList()
		exampleOwnerID := fakes.BuildFakeID()

		expectations := setupExpectations(
			mpm,
			func(db *database.MockDatabase) {
				db.InstrumentOwnershipDataManagerMock.On(testutils.GetMethodName(mpm.db.GetInstrumentOwnerships), testutils.ContextMatcher, exampleOwnerID, testutils.QueryFilterMatcher).Return(expected, nil)
			},
		)

		actual, cursor, err := mpm.ListInstrumentOwnerships(ctx, exampleOwnerID, nil)
		assert.NoError(t, err)
		assert.Empty(t, cursor)
		assert.Equal(t, expected.Data, actual)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestMealPlanningManager_CreateInstrumentOwnership(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		mpm := buildMealPlanManagerForTest(t)

		expected := fakes.BuildFakeInstrumentOwnership()
		fakeInput := fakes.BuildFakeInstrumentOwnershipCreationRequestInput()

		expectations := setupExpectations(
			mpm,
			func(db *database.MockDatabase) {
				db.InstrumentOwnershipDataManagerMock.On(testutils.GetMethodName(mpm.db.CreateInstrumentOwnership), testutils.ContextMatcher, testutils.MatchType[*types.InstrumentOwnershipDatabaseCreationInput]()).Return(expected, nil)
			},
			map[string][]string{
				events.InstrumentOwnershipCreated: {keys.HouseholdInstrumentOwnershipIDKey},
			},
		)

		actual, err := mpm.CreateInstrumentOwnership(ctx, fakeInput)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestMealPlanningManager_ReadInstrumentOwnership(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		mpm := buildMealPlanManagerForTest(t)

		ownerID := fakes.BuildFakeID()
		expected := fakes.BuildFakeInstrumentOwnership()

		expectations := setupExpectations(
			mpm,
			func(db *database.MockDatabase) {
				db.InstrumentOwnershipDataManagerMock.On(testutils.GetMethodName(mpm.db.GetInstrumentOwnership), testutils.ContextMatcher, expected.ID, ownerID).Return(expected, nil)
			},
		)

		actual, err := mpm.ReadInstrumentOwnership(ctx, ownerID, expected.ID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mock.AssertExpectationsForObjects(t, expectations...)
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
