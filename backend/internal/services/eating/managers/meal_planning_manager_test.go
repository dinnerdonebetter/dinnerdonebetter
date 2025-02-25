package managers

import (
	"testing"

	msgconfig "github.com/dinnerdonebetter/backend/internal/lib/messagequeue/config"
	mockpublishers "github.com/dinnerdonebetter/backend/internal/lib/messagequeue/mock"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/metrics"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
	textsearchcfg "github.com/dinnerdonebetter/backend/internal/lib/search/text/config"
	"github.com/dinnerdonebetter/backend/internal/lib/testutils"
	"github.com/dinnerdonebetter/backend/internal/services/eating/database"
	"github.com/dinnerdonebetter/backend/internal/services/eating/events"
	"github.com/dinnerdonebetter/backend/internal/services/eating/types"
	"github.com/dinnerdonebetter/backend/internal/services/eating/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func buildMealPlanManagerForTest(t *testing.T) *mealPlanningManager {
	t.Helper()

	queueCfg := &msgconfig.QueuesConfig{
		DataChangesTopicName: t.Name(),
	}

	mpp := &mockpublishers.PublisherProvider{}
	mpp.On("ProvidePublisher", queueCfg.DataChangesTopicName).Return(&mockpublishers.Publisher{}, nil)

	m, err := NewMealPlanningManager(
		t.Context(),
		logging.NewNoopLogger(),
		tracing.NewNoopTracerProvider(),
		database.NewMockDatabase(),
		queueCfg,
		mpp,
		&textsearchcfg.Config{},
		metrics.NewNoopMetricsProvider(),
	)
	require.NoError(t, err)

	mock.AssertExpectationsForObjects(t, mpp)

	return m.(*mealPlanningManager)
}

func TestMealPlanningManager_ListMeals(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		mpm := buildMealPlanManagerForTest(t)

		expected := fakes.BuildFakeMealsList()

		expectations := setupExpectationsForMealPlanningManager(
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

		expectations := setupExpectationsForMealPlanningManager(
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

		expectations := setupExpectationsForMealPlanningManager(
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

		expectations := setupExpectationsForMealPlanningManager(
			mpm,
			func(db *database.MockDatabase) {
				db.MealDataManagerMock.On(testutils.GetMethodName(mpm.db.SearchForMeals), testutils.ContextMatcher, exampleQuery, testutils.QueryFilterMatcher).Return(expected, nil)
			},
		)

		actual, err := mpm.SearchMeals(ctx, exampleQuery, true, nil)
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

		expectations := setupExpectationsForMealPlanningManager(
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

		expectations := setupExpectationsForMealPlanningManager(
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

		expectations := setupExpectationsForMealPlanningManager(
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

		expectations := setupExpectationsForMealPlanningManager(
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

		ctx := t.Context()
		mpm := buildMealPlanManagerForTest(t)

		exampleMealPlan := fakes.BuildFakeMealPlan()
		ownerID := fakes.BuildFakeID()
		exampleInput := fakes.BuildFakeMealPlanUpdateRequestInput()

		expectations := setupExpectationsForMealPlanningManager(
			mpm,
			func(db *database.MockDatabase) {
				db.MealPlanDataManagerMock.On(testutils.GetMethodName(mpm.db.GetMealPlan), testutils.ContextMatcher, exampleMealPlan.ID, ownerID).Return(exampleMealPlan, nil)
				db.MealPlanDataManagerMock.On(testutils.GetMethodName(mpm.db.UpdateMealPlan), testutils.ContextMatcher, testutils.MatchType[*types.MealPlan]()).Return(nil)
			},
			map[string][]string{
				events.MealPlanUpdated: {keys.MealPlanIDKey},
			},
		)

		assert.NoError(t, mpm.UpdateMealPlan(ctx, exampleMealPlan.ID, ownerID, exampleInput))

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestMealPlanningManager_ArchiveMealPlan(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		mpm := buildMealPlanManagerForTest(t)

		expected := fakes.BuildFakeMealPlan()

		expectations := setupExpectationsForMealPlanningManager(
			mpm,
			func(db *database.MockDatabase) {
				db.MealPlanDataManagerMock.On(testutils.GetMethodName(mpm.db.ArchiveMealPlan), testutils.ContextMatcher, expected.ID, expected.CreatedByUser).Return(nil)
			},
			map[string][]string{
				events.MealPlanArchived: {keys.MealPlanIDKey},
			},
		)

		err := mpm.ArchiveMealPlan(ctx, expected.ID, expected.CreatedByUser)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestMealPlanningManager_FinalizeMealPlan(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		mpm := buildMealPlanManagerForTest(t)

		expected := fakes.BuildFakeMealPlan()

		expectations := setupExpectationsForMealPlanningManager(
			mpm,
			func(db *database.MockDatabase) {
				db.MealPlanDataManagerMock.On(testutils.GetMethodName(mpm.db.AttemptToFinalizeMealPlan), testutils.ContextMatcher, expected.ID, expected.CreatedByUser).Return(true, nil)
			},
			map[string][]string{
				events.MealPlanFinalized: {keys.MealPlanIDKey},
			},
		)

		finalized, err := mpm.FinalizeMealPlan(ctx, expected.ID, expected.CreatedByUser)
		assert.True(t, finalized)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, expectations...)
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

		expectations := setupExpectationsForMealPlanningManager(
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

		expectations := setupExpectationsForMealPlanningManager(
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

		expectations := setupExpectationsForMealPlanningManager(
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

		ctx := t.Context()
		mpm := buildMealPlanManagerForTest(t)

		exampleMealPlanEvent := fakes.BuildFakeMealPlanEvent()
		exampleMealPlanID := fakes.BuildFakeID()
		exampleInput := fakes.BuildFakeMealPlanEventUpdateRequestInput()

		expectations := setupExpectationsForMealPlanningManager(
			mpm,
			func(db *database.MockDatabase) {
				db.MealPlanEventDataManagerMock.On(testutils.GetMethodName(mpm.db.GetMealPlanEvent), testutils.ContextMatcher, exampleMealPlanID, exampleMealPlanEvent.ID).Return(exampleMealPlanEvent, nil)
				db.MealPlanEventDataManagerMock.On(testutils.GetMethodName(mpm.db.UpdateMealPlanEvent), testutils.ContextMatcher, testutils.MatchType[*types.MealPlanEvent]()).Return(nil)
			},
			map[string][]string{
				events.MealPlanEventUpdated: {
					keys.MealPlanIDKey,
					keys.MealPlanEventIDKey,
				},
			},
		)

		assert.NoError(t, mpm.UpdateMealPlanEvent(ctx, exampleMealPlanID, exampleMealPlanEvent.ID, exampleInput))

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestMealPlanningManager_ArchiveMealPlanEvent(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		mpm := buildMealPlanManagerForTest(t)

		mealPlanID := fakes.BuildFakeID()
		expected := fakes.BuildFakeMealPlanEvent()

		expectations := setupExpectationsForMealPlanningManager(
			mpm,
			func(db *database.MockDatabase) {
				db.MealPlanEventDataManagerMock.On(testutils.GetMethodName(mpm.db.ArchiveMealPlanEvent), testutils.ContextMatcher, mealPlanID, expected.ID).Return(nil)
			},
			map[string][]string{
				events.MealPlanEventArchived: {
					keys.MealPlanIDKey,
					keys.MealPlanEventIDKey,
				},
			},
		)

		err := mpm.ArchiveMealPlanEvent(ctx, mealPlanID, expected.ID)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, expectations...)
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

		expectations := setupExpectationsForMealPlanningManager(
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

		expectations := setupExpectationsForMealPlanningManager(
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

		expectations := setupExpectationsForMealPlanningManager(
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

		ctx := t.Context()
		mpm := buildMealPlanManagerForTest(t)

		exampleMealPlanOption := fakes.BuildFakeMealPlanOption()
		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanEventID := fakes.BuildFakeID()
		exampleInput := fakes.BuildFakeMealPlanOptionUpdateRequestInput()

		expectations := setupExpectationsForMealPlanningManager(
			mpm,
			func(db *database.MockDatabase) {
				db.MealPlanOptionDataManagerMock.On(testutils.GetMethodName(mpm.db.GetMealPlanOption), testutils.ContextMatcher, exampleMealPlanID, exampleMealPlanEventID, exampleMealPlanOption.ID).Return(exampleMealPlanOption, nil)
				db.MealPlanOptionDataManagerMock.On(testutils.GetMethodName(mpm.db.UpdateMealPlanOption), testutils.ContextMatcher, testutils.MatchType[*types.MealPlanOption]()).Return(nil)
			},
			map[string][]string{
				events.MealPlanOptionUpdated: {
					keys.MealPlanIDKey,
					keys.MealPlanEventIDKey,
					keys.MealPlanOptionIDKey,
				},
			},
		)

		assert.NoError(t, mpm.UpdateMealPlanOption(ctx, exampleMealPlanID, exampleMealPlanEventID, exampleMealPlanOption.ID, exampleInput))

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestMealPlanningManager_ArchiveMealPlanOption(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		mpm := buildMealPlanManagerForTest(t)

		mealPlanID := fakes.BuildFakeID()
		mealPlanEventID := fakes.BuildFakeID()
		expected := fakes.BuildFakeMealPlanOption()

		expectations := setupExpectationsForMealPlanningManager(
			mpm,
			func(db *database.MockDatabase) {
				db.MealPlanOptionDataManagerMock.On(testutils.GetMethodName(mpm.db.ArchiveMealPlanOption), testutils.ContextMatcher, mealPlanID, mealPlanEventID, expected.ID).Return(nil)
			},
			map[string][]string{
				events.MealPlanOptionArchived: {
					keys.MealPlanIDKey,
					keys.MealPlanEventIDKey,
					keys.MealPlanOptionIDKey,
				},
			},
		)

		err := mpm.ArchiveMealPlanOption(ctx, mealPlanID, mealPlanEventID, expected.ID)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, expectations...)
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

		expectations := setupExpectationsForMealPlanningManager(
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

		expectations := setupExpectationsForMealPlanningManager(
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

		expectations := setupExpectationsForMealPlanningManager(
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

		ctx := t.Context()
		mpm := buildMealPlanManagerForTest(t)

		exampleMealPlanOptionVote := fakes.BuildFakeMealPlanOptionVote()
		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanOptionID := fakes.BuildFakeID()
		exampleMealPlanEventID := fakes.BuildFakeID()
		exampleInput := fakes.BuildFakeMealPlanOptionVoteUpdateRequestInput()

		expectations := setupExpectationsForMealPlanningManager(
			mpm,
			func(db *database.MockDatabase) {
				db.MealPlanOptionVoteDataManagerMock.On(testutils.GetMethodName(mpm.db.GetMealPlanOptionVote), testutils.ContextMatcher, exampleMealPlanID, exampleMealPlanEventID, exampleMealPlanOptionID, exampleMealPlanOptionVote.ID).Return(exampleMealPlanOptionVote, nil)
				db.MealPlanOptionVoteDataManagerMock.On(testutils.GetMethodName(mpm.db.UpdateMealPlanOptionVote), testutils.ContextMatcher, testutils.MatchType[*types.MealPlanOptionVote]()).Return(nil)
			},
			map[string][]string{
				events.MealPlanOptionVoteUpdated: {
					keys.MealPlanIDKey,
					keys.MealPlanEventIDKey,
					keys.MealPlanOptionIDKey,
					keys.MealPlanOptionVoteIDKey,
				},
			},
		)

		assert.NoError(t, mpm.UpdateMealPlanOptionVote(ctx, exampleMealPlanID, exampleMealPlanEventID, exampleMealPlanOptionID, exampleMealPlanOptionVote.ID, exampleInput))

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestMealPlanningManager_ArchiveMealPlanOptionVote(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		mpm := buildMealPlanManagerForTest(t)

		mealPlanID := fakes.BuildFakeID()
		mealPlanEventID := fakes.BuildFakeID()
		mealPlanOptionID := fakes.BuildFakeID()
		expected := fakes.BuildFakeMealPlanOptionVote()

		expectations := setupExpectationsForMealPlanningManager(
			mpm,
			func(db *database.MockDatabase) {
				db.MealPlanOptionVoteDataManagerMock.On(testutils.GetMethodName(mpm.db.ArchiveMealPlanOptionVote), testutils.ContextMatcher, mealPlanID, mealPlanEventID, mealPlanOptionID, expected.ID).Return(nil)
			},
			map[string][]string{
				events.MealPlanOptionVoteArchived: {
					keys.MealPlanIDKey,
					keys.MealPlanEventIDKey,
					keys.MealPlanOptionIDKey,
					keys.MealPlanOptionVoteIDKey,
				},
			},
		)

		err := mpm.ArchiveMealPlanOptionVote(ctx, mealPlanID, mealPlanEventID, mealPlanOptionID, expected.ID)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, expectations...)
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

		expectations := setupExpectationsForMealPlanningManager(
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

		expectations := setupExpectationsForMealPlanningManager(
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

		expectations := setupExpectationsForMealPlanningManager(
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

		ctx := t.Context()
		mpm := buildMealPlanManagerForTest(t)

		exampleInput := fakes.BuildFakeMealPlanTaskStatusChangeRequestInput()

		expectations := setupExpectationsForMealPlanningManager(
			mpm,
			func(db *database.MockDatabase) {
				db.MealPlanTaskDataManagerMock.On(testutils.GetMethodName(mpm.db.ChangeMealPlanTaskStatus), testutils.ContextMatcher, exampleInput).Return(nil)
			},
			map[string][]string{
				events.MealPlanTaskStatusChanged: {
					keys.MealPlanTaskIDKey,
				},
			},
		)

		assert.NoError(t, mpm.MealPlanTaskStatusChange(ctx, exampleInput))

		mock.AssertExpectationsForObjects(t, expectations...)
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

		expectations := setupExpectationsForMealPlanningManager(
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

		expectations := setupExpectationsForMealPlanningManager(
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

		expectations := setupExpectationsForMealPlanningManager(
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

		ctx := t.Context()
		mpm := buildMealPlanManagerForTest(t)

		exampleMealPlanGroceryListItem := fakes.BuildFakeMealPlanGroceryListItem()
		exampleMealPlanID := fakes.BuildFakeID()
		exampleInput := fakes.BuildFakeMealPlanGroceryListItemUpdateRequestInput()

		expectations := setupExpectationsForMealPlanningManager(
			mpm,
			func(db *database.MockDatabase) {
				db.MealPlanGroceryListItemDataManagerMock.On(testutils.GetMethodName(mpm.db.GetMealPlanGroceryListItem), testutils.ContextMatcher, exampleMealPlanID, exampleMealPlanGroceryListItem.ID).Return(exampleMealPlanGroceryListItem, nil)
				db.MealPlanGroceryListItemDataManagerMock.On(testutils.GetMethodName(mpm.db.UpdateMealPlanGroceryListItem), testutils.ContextMatcher, testutils.MatchType[*types.MealPlanGroceryListItem]()).Return(nil)
			},
			map[string][]string{
				events.MealPlanGroceryListItemUpdated: {
					keys.MealPlanIDKey,
					keys.MealPlanGroceryListItemIDKey,
				},
			},
		)

		assert.NoError(t, mpm.UpdateMealPlanGroceryListItem(ctx, exampleMealPlanID, exampleMealPlanGroceryListItem.ID, exampleInput))

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestMealPlanningManager_ArchiveMealPlanGroceryListItem(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		mpm := buildMealPlanManagerForTest(t)

		mealPlanID := fakes.BuildFakeID()
		expected := fakes.BuildFakeMealPlanGroceryListItem()

		expectations := setupExpectationsForMealPlanningManager(
			mpm,
			func(db *database.MockDatabase) {
				db.MealPlanGroceryListItemDataManagerMock.On(testutils.GetMethodName(mpm.db.ArchiveMealPlanGroceryListItem), testutils.ContextMatcher, mealPlanID, expected.ID).Return(nil)
			},
			map[string][]string{
				events.MealPlanGroceryListItemArchived: {
					keys.MealPlanIDKey,
					keys.MealPlanGroceryListItemIDKey,
				},
			},
		)

		err := mpm.ArchiveMealPlanGroceryListItem(ctx, mealPlanID, expected.ID)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, expectations...)
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

		expectations := setupExpectationsForMealPlanningManager(
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

		expectations := setupExpectationsForMealPlanningManager(
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

		ctx := t.Context()
		mpm := buildMealPlanManagerForTest(t)

		exampleIngredientPreference := fakes.BuildFakeIngredientPreference()
		ownerID := fakes.BuildFakeID()
		exampleInput := fakes.BuildFakeIngredientPreferenceUpdateRequestInput()

		expectations := setupExpectationsForMealPlanningManager(
			mpm,
			func(db *database.MockDatabase) {
				db.IngredientPreferenceDataManagerMock.On(testutils.GetMethodName(mpm.db.GetIngredientPreference), testutils.ContextMatcher, exampleIngredientPreference.ID, ownerID).Return(exampleIngredientPreference, nil)
				db.IngredientPreferenceDataManagerMock.On(testutils.GetMethodName(mpm.db.UpdateIngredientPreference), testutils.ContextMatcher, testutils.MatchType[*types.IngredientPreference]()).Return(nil)
			},
			map[string][]string{
				events.IngredientPreferenceUpdated: {
					keys.IngredientPreferenceIDKey,
				},
			},
		)

		assert.NoError(t, mpm.UpdateIngredientPreference(ctx, exampleIngredientPreference.ID, ownerID, exampleInput))

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestMealPlanningManager_ArchiveIngredientPreference(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		mpm := buildMealPlanManagerForTest(t)

		ownershipID := fakes.BuildFakeID()
		expected := fakes.BuildFakeIngredientPreference()

		expectations := setupExpectationsForMealPlanningManager(
			mpm,
			func(db *database.MockDatabase) {
				db.IngredientPreferenceDataManagerMock.On(testutils.GetMethodName(mpm.db.ArchiveIngredientPreference), testutils.ContextMatcher, ownershipID, expected.ID).Return(nil)
			},
			map[string][]string{
				events.IngredientPreferenceArchived: {
					keys.IngredientPreferenceIDKey,
				},
			},
		)

		err := mpm.ArchiveIngredientPreference(ctx, ownershipID, expected.ID)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, expectations...)
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

		expectations := setupExpectationsForMealPlanningManager(
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

		expectations := setupExpectationsForMealPlanningManager(
			mpm,
			func(db *database.MockDatabase) {
				db.InstrumentOwnershipDataManagerMock.On(testutils.GetMethodName(mpm.db.CreateInstrumentOwnership), testutils.ContextMatcher, testutils.MatchType[*types.InstrumentOwnershipDatabaseCreationInput]()).Return(expected, nil)
			},
			map[string][]string{
				events.InstrumentOwnershipCreated: {keys.InstrumentOwnershipIDKey},
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

		expectations := setupExpectationsForMealPlanningManager(
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

		ctx := t.Context()
		mpm := buildMealPlanManagerForTest(t)

		exampleInstrumentOwnership := fakes.BuildFakeInstrumentOwnership()
		ownerID := fakes.BuildFakeID()
		exampleInput := fakes.BuildFakeInstrumentOwnershipUpdateRequestInput()

		expectations := setupExpectationsForMealPlanningManager(
			mpm,
			func(db *database.MockDatabase) {
				db.InstrumentOwnershipDataManagerMock.On(testutils.GetMethodName(mpm.db.GetInstrumentOwnership), testutils.ContextMatcher, exampleInstrumentOwnership.ID, ownerID).Return(exampleInstrumentOwnership, nil)
				db.InstrumentOwnershipDataManagerMock.On(testutils.GetMethodName(mpm.db.UpdateInstrumentOwnership), testutils.ContextMatcher, testutils.MatchType[*types.InstrumentOwnership]()).Return(nil)
			},
			map[string][]string{
				events.InstrumentOwnershipUpdated: {
					keys.InstrumentOwnershipIDKey,
				},
			},
		)

		assert.NoError(t, mpm.UpdateInstrumentOwnership(ctx, exampleInstrumentOwnership.ID, ownerID, exampleInput))

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestMealPlanningManager_ArchiveInstrumentOwnership(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		mpm := buildMealPlanManagerForTest(t)

		ownershipID := fakes.BuildFakeID()
		expected := fakes.BuildFakeInstrumentOwnership()

		expectations := setupExpectationsForMealPlanningManager(
			mpm,
			func(db *database.MockDatabase) {
				db.InstrumentOwnershipDataManagerMock.On(testutils.GetMethodName(mpm.db.ArchiveInstrumentOwnership), testutils.ContextMatcher, expected.ID, ownershipID).Return(nil)
			},
			map[string][]string{
				events.InstrumentOwnershipArchived: {
					keys.InstrumentOwnershipIDKey,
				},
			},
		)

		err := mpm.ArchiveInstrumentOwnership(ctx, ownershipID, expected.ID)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}
