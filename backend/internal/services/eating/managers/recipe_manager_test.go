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
	"github.com/dinnerdonebetter/backend/internal/services/eating/businesslogic/recipeanalysis"
	"github.com/dinnerdonebetter/backend/internal/services/eating/database"
	"github.com/dinnerdonebetter/backend/internal/services/eating/events"
	"github.com/dinnerdonebetter/backend/internal/services/eating/types"
	"github.com/dinnerdonebetter/backend/internal/services/eating/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func buildRecipeManagerForTest(t *testing.T) *recipeManager {
	t.Helper()

	queueCfg := &msgconfig.QueuesConfig{
		DataChangesTopicName: t.Name(),
	}

	mpp := &mockpublishers.PublisherProvider{}
	mpp.On("ProvidePublisher", queueCfg.DataChangesTopicName).Return(&mockpublishers.Publisher{}, nil)

	m, err := NewRecipeManager(
		t.Context(),
		logging.NewNoopLogger(),
		tracing.NewNoopTracerProvider(),
		database.NewMockDatabase(),
		queueCfg,
		mpp,
		&recipeanalysis.MockRecipeAnalyzer{},
		&textsearchcfg.Config{},
		metrics.NewNoopMetricsProvider(),
	)
	require.NoError(t, err)

	mock.AssertExpectationsForObjects(t, mpp)

	return m.(*recipeManager)
}

func setupExpectationsForRecipeManager(
	manager *recipeManager,
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

func TestRecipeManager_ListRecipes(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		rm := buildRecipeManagerForTest(t)

		expected := fakes.BuildFakeRecipesList()

		expectations := setupExpectationsForRecipeManager(
			rm,
			func(db *database.MockDatabase) {
				db.RecipeDataManagerMock.On(testutils.GetMethodName(rm.db.GetRecipes), testutils.ContextMatcher, testutils.QueryFilterMatcher).Return(expected, nil)
			},
		)

		actual, cursor, err := rm.ListRecipes(ctx, nil)
		assert.NoError(t, err)
		assert.Empty(t, cursor)
		assert.Equal(t, expected.Data, actual)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestRecipeManager_CreateRecipe(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		rm := buildRecipeManagerForTest(t)

		expected := fakes.BuildFakeRecipe()
		fakeInput := fakes.BuildFakeRecipeCreationRequestInput()

		expectations := setupExpectationsForRecipeManager(
			rm,
			func(db *database.MockDatabase) {
				db.RecipeDataManagerMock.On(testutils.GetMethodName(rm.db.CreateRecipe), testutils.ContextMatcher, testutils.MatchType[*types.RecipeDatabaseCreationInput]()).Return(expected, nil)
			},
			map[string][]string{
				events.RecipeCreated: {keys.RecipeIDKey},
			},
		)

		actual, err := rm.CreateRecipe(ctx, fakeInput)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestRecipeManager_ReadRecipe(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		rm := buildRecipeManagerForTest(t)

		expected := fakes.BuildFakeRecipe()

		expectations := setupExpectationsForRecipeManager(
			rm,
			func(db *database.MockDatabase) {
				db.RecipeDataManagerMock.On(testutils.GetMethodName(rm.db.GetRecipe), testutils.ContextMatcher, expected.ID).Return(expected, nil)
			},
		)

		actual, err := rm.ReadRecipe(ctx, expected.ID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestRecipeManager_SearchRecipes(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		rm := buildRecipeManagerForTest(t)

		expected := fakes.BuildFakeRecipesList()
		exampleQuery := fakes.BuildFakeID()

		expectations := setupExpectationsForRecipeManager(
			rm,
			func(db *database.MockDatabase) {
				db.RecipeDataManagerMock.On(testutils.GetMethodName(rm.db.SearchForRecipes), testutils.ContextMatcher, exampleQuery, testutils.QueryFilterMatcher).Return(expected, nil)
			},
		)

		actual, cursor, err := rm.SearchRecipes(ctx, exampleQuery, true, nil)
		assert.NoError(t, err)
		assert.Empty(t, cursor)
		assert.Equal(t, expected.Data, actual)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestRecipeManager_UpdateRecipe(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		rm := buildRecipeManagerForTest(t)

		exampleRecipe := fakes.BuildFakeRecipe()
		exampleInput := fakes.BuildFakeRecipeUpdateRequestInput()

		expectations := setupExpectationsForRecipeManager(
			rm,
			func(db *database.MockDatabase) {
				db.RecipeDataManagerMock.On(testutils.GetMethodName(rm.db.GetRecipe), testutils.ContextMatcher, exampleRecipe.ID).Return(exampleRecipe, nil)
				db.RecipeDataManagerMock.On(testutils.GetMethodName(rm.db.UpdateRecipe), testutils.ContextMatcher, testutils.MatchType[*types.Recipe]()).Return(nil)
			},
			map[string][]string{
				events.RecipeUpdated: {keys.RecipeIDKey},
			},
		)

		assert.NoError(t, rm.UpdateRecipe(ctx, exampleRecipe.ID, exampleInput))

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestRecipeManager_ArchiveRecipe(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		rm := buildRecipeManagerForTest(t)

		expected := fakes.BuildFakeRecipe()
		exampleOwnerID := fakes.BuildFakeID()

		expectations := setupExpectationsForRecipeManager(
			rm,
			func(db *database.MockDatabase) {
				db.RecipeDataManagerMock.On(testutils.GetMethodName(rm.db.ArchiveRecipe), testutils.ContextMatcher, expected.ID, exampleOwnerID).Return(nil)
			},
			map[string][]string{
				events.RecipeArchived: {keys.RecipeIDKey},
			},
		)

		assert.NoError(t, rm.ArchiveRecipe(ctx, expected.ID, exampleOwnerID))

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestRecipeManager_RecipeEstimatedPrepSteps(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		t.SkipNow()
	})
}

func TestRecipeManager_RecipeImageUpload(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		t.SkipNow()
	})
}

func TestRecipeManager_RecipeMermaid(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		t.SkipNow()
	})
}

func TestRecipeManager_CloneRecipe(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		t.SkipNow()
	})
}

func TestRecipeManager_ListRecipeSteps(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		rm := buildRecipeManagerForTest(t)

		expected := fakes.BuildFakeRecipeStepsList()
		exampleRecipeID := fakes.BuildFakeID()

		expectations := setupExpectationsForRecipeManager(
			rm,
			func(db *database.MockDatabase) {
				db.RecipeStepDataManagerMock.On(testutils.GetMethodName(rm.db.GetRecipeSteps), testutils.ContextMatcher, exampleRecipeID, testutils.QueryFilterMatcher).Return(expected, nil)
			},
		)

		actual, cursor, err := rm.ListRecipeSteps(ctx, exampleRecipeID, nil)
		assert.NoError(t, err)
		assert.Empty(t, cursor)
		assert.Equal(t, expected.Data, actual)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestRecipeManager_CreateRecipeStep(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		t.SkipNow()
	})
}

func TestRecipeManager_ReadRecipeStep(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		t.SkipNow()
	})
}

func TestRecipeManager_UpdateRecipeStep(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		t.SkipNow()
	})
}

func TestRecipeManager_ArchiveRecipeStep(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		t.SkipNow()
	})
}

func TestRecipeManager_RecipeStepImageUpload(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		t.SkipNow()
	})
}

func TestRecipeManager_ListRecipeStepProducts(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		rm := buildRecipeManagerForTest(t)

		expected := fakes.BuildFakeRecipeStepProductsList()
		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()

		expectations := setupExpectationsForRecipeManager(
			rm,
			func(db *database.MockDatabase) {
				db.RecipeStepProductDataManagerMock.On(testutils.GetMethodName(rm.db.GetRecipeStepProducts), testutils.ContextMatcher, exampleRecipeID, exampleRecipeStepID, testutils.QueryFilterMatcher).Return(expected, nil)
			},
		)

		actual, cursor, err := rm.ListRecipeStepProducts(ctx, exampleRecipeID, exampleRecipeStepID, nil)
		assert.NoError(t, err)
		assert.Empty(t, cursor)
		assert.Equal(t, expected.Data, actual)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestRecipeManager_CreateRecipeStepProduct(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		t.SkipNow()
	})
}

func TestRecipeManager_ReadRecipeStepProduct(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		t.SkipNow()
	})
}

func TestRecipeManager_UpdateRecipeStepProduct(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		t.SkipNow()
	})
}

func TestRecipeManager_ArchiveRecipeStepProduct(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		t.SkipNow()
	})
}

func TestRecipeManager_ListRecipeStepInstruments(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		rm := buildRecipeManagerForTest(t)

		expected := fakes.BuildFakeRecipeStepInstrumentsList()
		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()

		expectations := setupExpectationsForRecipeManager(
			rm,
			func(db *database.MockDatabase) {
				db.RecipeStepInstrumentDataManagerMock.On(testutils.GetMethodName(rm.db.GetRecipeStepInstruments), testutils.ContextMatcher, exampleRecipeID, exampleRecipeStepID, testutils.QueryFilterMatcher).Return(expected, nil)
			},
		)

		actual, cursor, err := rm.ListRecipeStepInstruments(ctx, exampleRecipeID, exampleRecipeStepID, nil)
		assert.NoError(t, err)
		assert.Empty(t, cursor)
		assert.Equal(t, expected.Data, actual)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestRecipeManager_CreateRecipeStepInstrument(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		t.SkipNow()
	})
}

func TestRecipeManager_ReadRecipeStepInstrument(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		t.SkipNow()
	})
}

func TestRecipeManager_UpdateRecipeStepInstrument(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		t.SkipNow()
	})
}

func TestRecipeManager_ArchiveRecipeStepInstrument(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		t.SkipNow()
	})
}

func TestRecipeManager_ListRecipeStepIngredients(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		rm := buildRecipeManagerForTest(t)

		expected := fakes.BuildFakeRecipeStepIngredientsList()
		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()

		expectations := setupExpectationsForRecipeManager(
			rm,
			func(db *database.MockDatabase) {
				db.RecipeStepIngredientDataManagerMock.On(testutils.GetMethodName(rm.db.GetRecipeStepIngredients), testutils.ContextMatcher, exampleRecipeID, exampleRecipeStepID, testutils.QueryFilterMatcher).Return(expected, nil)
			},
		)

		actual, cursor, err := rm.ListRecipeStepIngredients(ctx, exampleRecipeID, exampleRecipeStepID, nil)
		assert.NoError(t, err)
		assert.Empty(t, cursor)
		assert.Equal(t, expected.Data, actual)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestRecipeManager_CreateRecipeStepIngredient(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		t.SkipNow()
	})
}

func TestRecipeManager_ReadRecipeStepIngredient(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		t.SkipNow()
	})
}

func TestRecipeManager_UpdateRecipeStepIngredient(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		t.SkipNow()
	})
}

func TestRecipeManager_ArchiveRecipeStepIngredient(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		t.SkipNow()
	})
}

func TestRecipeManager_ListRecipePrepTask(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		rm := buildRecipeManagerForTest(t)

		expected := fakes.BuildFakeRecipePrepTasksList()
		exampleRecipeID := fakes.BuildFakeID()

		expectations := setupExpectationsForRecipeManager(
			rm,
			func(db *database.MockDatabase) {
				db.RecipePrepTaskDataManagerMock.On(testutils.GetMethodName(rm.db.GetRecipePrepTasksForRecipe), testutils.ContextMatcher, exampleRecipeID).Return(expected.Data, nil)
			},
		)

		actual, cursor, err := rm.ListRecipePrepTask(ctx, exampleRecipeID, nil)
		assert.NoError(t, err)
		assert.Empty(t, cursor)
		assert.Equal(t, expected.Data, actual)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestRecipeManager_CreateRecipePrepTask(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		t.SkipNow()
	})
}

func TestRecipeManager_ReadRecipePrepTask(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		t.SkipNow()
	})
}

func TestRecipeManager_UpdateRecipePrepTask(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		t.SkipNow()
	})
}

func TestRecipeManager_ArchiveRecipePrepTask(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		t.SkipNow()
	})
}

func TestRecipeManager_ListRecipeStepCompletionConditions(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		rm := buildRecipeManagerForTest(t)

		expected := fakes.BuildFakeRecipeStepCompletionConditionsList()
		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()

		expectations := setupExpectationsForRecipeManager(
			rm,
			func(db *database.MockDatabase) {
				db.RecipeStepCompletionConditionDataManagerMock.On(testutils.GetMethodName(rm.db.GetRecipeStepCompletionConditions), testutils.ContextMatcher, exampleRecipeID, exampleRecipeStepID, testutils.QueryFilterMatcher).Return(expected, nil)
			},
		)

		actual, cursor, err := rm.ListRecipeStepCompletionConditions(ctx, exampleRecipeID, exampleRecipeStepID, nil)
		assert.NoError(t, err)
		assert.Empty(t, cursor)
		assert.Equal(t, expected.Data, actual)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestRecipeManager_CreateRecipeStepCompletionCondition(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		t.SkipNow()
	})
}

func TestRecipeManager_ReadRecipeStepCompletionCondition(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		t.SkipNow()
	})
}

func TestRecipeManager_UpdateRecipeStepCompletionCondition(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		t.SkipNow()
	})
}

func TestRecipeManager_ArchiveRecipeStepCompletionCondition(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		t.SkipNow()
	})
}

func TestRecipeManager_ListRecipeStepVessels(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		rm := buildRecipeManagerForTest(t)

		expected := fakes.BuildFakeRecipeStepVesselsList()
		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()

		expectations := setupExpectationsForRecipeManager(
			rm,
			func(db *database.MockDatabase) {
				db.RecipeStepVesselDataManagerMock.On(testutils.GetMethodName(rm.db.GetRecipeStepVessels), testutils.ContextMatcher, exampleRecipeID, exampleRecipeStepID, testutils.QueryFilterMatcher).Return(expected, nil)
			},
		)

		actual, cursor, err := rm.ListRecipeStepVessels(ctx, exampleRecipeID, exampleRecipeStepID, nil)
		assert.NoError(t, err)
		assert.Empty(t, cursor)
		assert.Equal(t, expected.Data, actual)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestRecipeManager_CreateRecipeStepVessel(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		t.SkipNow()
	})
}

func TestRecipeManager_ReadRecipeStepVessel(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		t.SkipNow()
	})
}

func TestRecipeManager_UpdateRecipeStepVessel(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		t.SkipNow()
	})
}

func TestRecipeManager_ArchiveRecipeStepVessel(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		t.SkipNow()
	})
}

func TestRecipeManager_ListRecipeRatings(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		rm := buildRecipeManagerForTest(t)

		expected := fakes.BuildFakeRecipeRatingsList()
		exampleRecipeID := fakes.BuildFakeID()

		expectations := setupExpectationsForRecipeManager(
			rm,
			func(db *database.MockDatabase) {
				db.RecipeRatingDataManagerMock.On(testutils.GetMethodName(rm.db.GetRecipeRatingsForRecipe), testutils.ContextMatcher, exampleRecipeID, testutils.QueryFilterMatcher).Return(expected, nil)
			},
		)

		actual, cursor, err := rm.ListRecipeRatings(ctx, exampleRecipeID, nil)
		assert.NoError(t, err)
		assert.Empty(t, cursor)
		assert.Equal(t, expected.Data, actual)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestRecipeManager_ReadRecipeRating(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		t.SkipNow()
	})
}

func TestRecipeManager_CreateRecipeRating(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		t.SkipNow()
	})
}

func TestRecipeManager_UpdateRecipeRating(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		t.SkipNow()
	})
}

func TestRecipeManager_ArchiveRecipeRating(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		t.SkipNow()
	})
}
