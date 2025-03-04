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

func setupExpectationsForRecipeManagerWithAnalyzer(
	manager *recipeManager,
	dbSetupFunc func(db *database.MockDatabase),
	analyzerSetupFunc func(analyzer *recipeanalysis.MockRecipeAnalyzer),
	eventTypeMaps ...map[string][]string,
) []any {
	expectations := setupExpectationsForRecipeManager(manager, dbSetupFunc, eventTypeMaps...)

	ra := &recipeanalysis.MockRecipeAnalyzer{}
	if analyzerSetupFunc != nil {
		analyzerSetupFunc(ra)
	}
	manager.recipeAnalyzer = ra

	return expectations
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
				events.RecipeCreated: {
					keys.RecipeIDKey,
				},
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
				events.RecipeUpdated: {
					keys.RecipeIDKey,
				},
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
				events.RecipeArchived: {
					keys.RecipeIDKey,
				},
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

		ctx := t.Context()
		rm := buildRecipeManagerForTest(t)

		exampleRecipe := fakes.BuildFakeRecipe()
		expectedResults := fakes.BuildFakeMealPlanTaskDatabaseCreationInputs()

		expectations := setupExpectationsForRecipeManagerWithAnalyzer(
			rm,
			func(db *database.MockDatabase) {
				db.RecipeDataManagerMock.On(testutils.GetMethodName(rm.db.GetRecipe), testutils.ContextMatcher, exampleRecipe.ID).Return(exampleRecipe, nil)
			},
			func(analyzer *recipeanalysis.MockRecipeAnalyzer) {
				analyzer.On("GenerateMealPlanTasksForRecipe", testutils.ContextMatcher, "", testutils.MatchType[*types.Recipe]()).Return(expectedResults, nil)
			},
		)

		results, err := rm.RecipeEstimatedPrepSteps(ctx, exampleRecipe.ID)
		assert.NoError(t, err)
		// TODO: find a better/cleaner way of asserting successful conversion
		assert.Equal(t, len(results), len(expectedResults))

		mock.AssertExpectationsForObjects(t, expectations...)
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

		ctx := t.Context()
		rm := buildRecipeManagerForTest(t)

		exampleRecipe := fakes.BuildFakeRecipe()
		expectedResult := t.Name()

		expectations := setupExpectationsForRecipeManagerWithAnalyzer(
			rm,
			func(db *database.MockDatabase) {
				db.RecipeDataManagerMock.On(testutils.GetMethodName(rm.db.GetRecipe), testutils.ContextMatcher, exampleRecipe.ID).Return(exampleRecipe, nil)
			},
			func(analyzer *recipeanalysis.MockRecipeAnalyzer) {
				analyzer.On("RenderMermaidDiagramForRecipe", testutils.ContextMatcher, testutils.MatchType[*types.Recipe]()).Return(expectedResult, nil)
			},
		)

		result, err := rm.RecipeMermaid(ctx, exampleRecipe.ID)
		assert.NoError(t, err)
		assert.Equal(t, result, expectedResult)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestRecipeManager_CloneRecipe(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		rm := buildRecipeManagerForTest(t)

		expected := fakes.BuildFakeRecipe()
		cloned := fakes.BuildFakeRecipe()
		exampleOwnerID := fakes.BuildFakeID()

		expectations := setupExpectationsForRecipeManager(
			rm,
			func(db *database.MockDatabase) {
				db.RecipeDataManagerMock.On(testutils.GetMethodName(rm.db.GetRecipe), testutils.ContextMatcher, expected.ID).Return(expected, nil)
				db.RecipeDataManagerMock.On(testutils.GetMethodName(rm.db.CreateRecipe), testutils.ContextMatcher, testutils.MatchType[*types.RecipeDatabaseCreationInput]()).Return(cloned, nil)
			},
			map[string][]string{
				events.RecipeCloned: {
					keys.RecipeIDKey,
				},
			},
		)

		actual, err := rm.CloneRecipe(ctx, expected.ID, exampleOwnerID)
		assert.NoError(t, err)
		assert.Equal(t, cloned, actual)

		mock.AssertExpectationsForObjects(t, expectations...)
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

		ctx := t.Context()
		rm := buildRecipeManagerForTest(t)

		exampleRecipeID := fakes.BuildFakeID()
		expected := fakes.BuildFakeRecipeStep()
		fakeInput := fakes.BuildFakeRecipeStepCreationRequestInput()

		expectations := setupExpectationsForRecipeManager(
			rm,
			func(db *database.MockDatabase) {
				db.RecipeStepDataManagerMock.On(testutils.GetMethodName(rm.db.CreateRecipeStep), testutils.ContextMatcher, testutils.MatchType[*types.RecipeStepDatabaseCreationInput]()).Return(expected, nil)
			},
			map[string][]string{
				events.RecipeStepCreated: {
					keys.RecipeIDKey,
					keys.RecipeStepIDKey,
				},
			},
		)

		actual, err := rm.CreateRecipeStep(ctx, exampleRecipeID, fakeInput)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestRecipeManager_ReadRecipeStep(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		rm := buildRecipeManagerForTest(t)

		exampleRecipeID := fakes.BuildFakeID()
		expected := fakes.BuildFakeRecipeStep()

		expectations := setupExpectationsForRecipeManager(
			rm,
			func(db *database.MockDatabase) {
				db.RecipeStepDataManagerMock.On(testutils.GetMethodName(rm.db.GetRecipeStep), testutils.ContextMatcher, exampleRecipeID, expected.ID).Return(expected, nil)
			},
		)

		actual, err := rm.ReadRecipeStep(ctx, exampleRecipeID, expected.ID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestRecipeManager_UpdateRecipeStep(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		rm := buildRecipeManagerForTest(t)

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStep := fakes.BuildFakeRecipeStep()
		exampleInput := fakes.BuildFakeRecipeStepUpdateRequestInput()

		expectations := setupExpectationsForRecipeManager(
			rm,
			func(db *database.MockDatabase) {
				db.RecipeStepDataManagerMock.On(testutils.GetMethodName(rm.db.GetRecipeStep), testutils.ContextMatcher, exampleRecipeID, exampleRecipeStep.ID).Return(exampleRecipeStep, nil)
				db.RecipeStepDataManagerMock.On(testutils.GetMethodName(rm.db.UpdateRecipeStep), testutils.ContextMatcher, testutils.MatchType[*types.RecipeStep]()).Return(nil)
			},
			map[string][]string{
				events.RecipeStepUpdated: {
					keys.RecipeIDKey,
					keys.RecipeStepIDKey,
				},
			},
		)

		assert.NoError(t, rm.UpdateRecipeStep(ctx, exampleRecipeID, exampleRecipeStep.ID, exampleInput))

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestRecipeManager_ArchiveRecipeStep(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		rm := buildRecipeManagerForTest(t)

		exampleRecipeID := fakes.BuildFakeID()
		expected := fakes.BuildFakeRecipeStep()

		expectations := setupExpectationsForRecipeManager(
			rm,
			func(db *database.MockDatabase) {
				db.RecipeStepDataManagerMock.On(testutils.GetMethodName(rm.db.ArchiveRecipeStep), testutils.ContextMatcher, exampleRecipeID, expected.ID).Return(nil)
			},
			map[string][]string{
				events.RecipeStepArchived: {
					keys.RecipeIDKey,
					keys.RecipeStepIDKey,
				},
			},
		)

		assert.NoError(t, rm.ArchiveRecipeStep(ctx, exampleRecipeID, expected.ID))

		mock.AssertExpectationsForObjects(t, expectations...)
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

		ctx := t.Context()
		rm := buildRecipeManagerForTest(t)

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()
		expected := fakes.BuildFakeRecipeStepProduct()
		fakeInput := fakes.BuildFakeRecipeStepProductCreationRequestInput()

		expectations := setupExpectationsForRecipeManager(
			rm,
			func(db *database.MockDatabase) {
				db.RecipeStepProductDataManagerMock.On(testutils.GetMethodName(rm.db.CreateRecipeStepProduct), testutils.ContextMatcher, testutils.MatchType[*types.RecipeStepProductDatabaseCreationInput]()).Return(expected, nil)
			},
			map[string][]string{
				events.RecipeStepProductCreated: {
					keys.RecipeIDKey,
					keys.RecipeStepIDKey,
					keys.RecipeStepProductIDKey,
				},
			},
		)

		actual, err := rm.CreateRecipeStepProduct(ctx, exampleRecipeID, exampleRecipeStepID, fakeInput)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestRecipeManager_ReadRecipeStepProduct(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		rm := buildRecipeManagerForTest(t)

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()
		expected := fakes.BuildFakeRecipeStepProduct()

		expectations := setupExpectationsForRecipeManager(
			rm,
			func(db *database.MockDatabase) {
				db.RecipeStepProductDataManagerMock.On(testutils.GetMethodName(rm.db.GetRecipeStepProduct), testutils.ContextMatcher, exampleRecipeID, exampleRecipeStepID, expected.ID).Return(expected, nil)
			},
		)

		actual, err := rm.ReadRecipeStepProduct(ctx, exampleRecipeID, exampleRecipeStepID, expected.ID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestRecipeManager_UpdateRecipeStepProduct(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		rm := buildRecipeManagerForTest(t)

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepProduct := fakes.BuildFakeRecipeStepProduct()
		exampleInput := fakes.BuildFakeRecipeStepProductUpdateRequestInput()

		expectations := setupExpectationsForRecipeManager(
			rm,
			func(db *database.MockDatabase) {
				db.RecipeStepProductDataManagerMock.On(testutils.GetMethodName(rm.db.GetRecipeStepProduct), testutils.ContextMatcher, exampleRecipeID, exampleRecipeStepID, exampleRecipeStepProduct.ID).Return(exampleRecipeStepProduct, nil)
				db.RecipeStepProductDataManagerMock.On(testutils.GetMethodName(rm.db.UpdateRecipeStepProduct), testutils.ContextMatcher, testutils.MatchType[*types.RecipeStepProduct]()).Return(nil)
			},
			map[string][]string{
				events.RecipeStepProductUpdated: {
					keys.RecipeIDKey,
					keys.RecipeStepIDKey,
					keys.RecipeStepProductIDKey,
				},
			},
		)

		assert.NoError(t, rm.UpdateRecipeStepProduct(ctx, exampleRecipeID, exampleRecipeStepID, exampleRecipeStepProduct.ID, exampleInput))

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestRecipeManager_ArchiveRecipeStepProduct(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		rm := buildRecipeManagerForTest(t)

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()
		expected := fakes.BuildFakeRecipeStepProduct()

		expectations := setupExpectationsForRecipeManager(
			rm,
			func(db *database.MockDatabase) {
				db.RecipeStepProductDataManagerMock.On(testutils.GetMethodName(rm.db.ArchiveRecipeStepProduct), testutils.ContextMatcher, exampleRecipeStepID, expected.ID).Return(nil)
			},
			map[string][]string{
				events.RecipeStepProductArchived: {
					keys.RecipeIDKey,
					keys.RecipeStepIDKey,
					keys.RecipeStepProductIDKey,
				},
			},
		)

		assert.NoError(t, rm.ArchiveRecipeStepProduct(ctx, exampleRecipeID, exampleRecipeStepID, expected.ID))

		mock.AssertExpectationsForObjects(t, expectations...)
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

		ctx := t.Context()
		rm := buildRecipeManagerForTest(t)

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()
		expected := fakes.BuildFakeRecipeStepInstrument()
		fakeInput := fakes.BuildFakeRecipeStepInstrumentCreationRequestInput()

		expectations := setupExpectationsForRecipeManager(
			rm,
			func(db *database.MockDatabase) {
				db.RecipeStepInstrumentDataManagerMock.On(testutils.GetMethodName(rm.db.CreateRecipeStepInstrument), testutils.ContextMatcher, testutils.MatchType[*types.RecipeStepInstrumentDatabaseCreationInput]()).Return(expected, nil)
			},
			map[string][]string{
				events.RecipeStepInstrumentCreated: {
					keys.RecipeIDKey,
					keys.RecipeStepIDKey,
					keys.RecipeStepInstrumentIDKey,
				},
			},
		)

		actual, err := rm.CreateRecipeStepInstrument(ctx, exampleRecipeID, exampleRecipeStepID, fakeInput)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestRecipeManager_ReadRecipeStepInstrument(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		rm := buildRecipeManagerForTest(t)

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()
		expected := fakes.BuildFakeRecipeStepInstrument()

		expectations := setupExpectationsForRecipeManager(
			rm,
			func(db *database.MockDatabase) {
				db.RecipeStepInstrumentDataManagerMock.On(testutils.GetMethodName(rm.db.GetRecipeStepInstrument), testutils.ContextMatcher, exampleRecipeID, exampleRecipeStepID, expected.ID).Return(expected, nil)
			},
		)

		actual, err := rm.ReadRecipeStepInstrument(ctx, exampleRecipeID, exampleRecipeStepID, expected.ID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestRecipeManager_UpdateRecipeStepInstrument(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		rm := buildRecipeManagerForTest(t)

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepInstrument := fakes.BuildFakeRecipeStepInstrument()
		exampleInput := fakes.BuildFakeRecipeStepInstrumentUpdateRequestInput()

		expectations := setupExpectationsForRecipeManager(
			rm,
			func(db *database.MockDatabase) {
				db.RecipeStepInstrumentDataManagerMock.On(testutils.GetMethodName(rm.db.GetRecipeStepInstrument), testutils.ContextMatcher, exampleRecipeID, exampleRecipeStepID, exampleRecipeStepInstrument.ID).Return(exampleRecipeStepInstrument, nil)
				db.RecipeStepInstrumentDataManagerMock.On(testutils.GetMethodName(rm.db.UpdateRecipeStepInstrument), testutils.ContextMatcher, testutils.MatchType[*types.RecipeStepInstrument]()).Return(nil)
			},
			map[string][]string{
				events.RecipeStepInstrumentUpdated: {
					keys.RecipeIDKey,
					keys.RecipeStepIDKey,
					keys.RecipeStepInstrumentIDKey,
				},
			},
		)

		assert.NoError(t, rm.UpdateRecipeStepInstrument(ctx, exampleRecipeID, exampleRecipeStepID, exampleRecipeStepInstrument.ID, exampleInput))

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestRecipeManager_ArchiveRecipeStepInstrument(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		rm := buildRecipeManagerForTest(t)

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()
		expected := fakes.BuildFakeRecipeStepInstrument()

		expectations := setupExpectationsForRecipeManager(
			rm,
			func(db *database.MockDatabase) {
				db.RecipeStepInstrumentDataManagerMock.On(testutils.GetMethodName(rm.db.ArchiveRecipeStepInstrument), testutils.ContextMatcher, exampleRecipeStepID, expected.ID).Return(nil)
			},
			map[string][]string{
				events.RecipeStepInstrumentArchived: {
					keys.RecipeIDKey,
					keys.RecipeStepIDKey,
					keys.RecipeStepInstrumentIDKey,
				},
			},
		)

		assert.NoError(t, rm.ArchiveRecipeStepInstrument(ctx, exampleRecipeID, exampleRecipeStepID, expected.ID))

		mock.AssertExpectationsForObjects(t, expectations...)
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

		ctx := t.Context()
		rm := buildRecipeManagerForTest(t)

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()
		expected := fakes.BuildFakeRecipeStepIngredient()
		fakeInput := fakes.BuildFakeRecipeStepIngredientCreationRequestInput()

		expectations := setupExpectationsForRecipeManager(
			rm,
			func(db *database.MockDatabase) {
				db.RecipeStepIngredientDataManagerMock.On(testutils.GetMethodName(rm.db.CreateRecipeStepIngredient), testutils.ContextMatcher, testutils.MatchType[*types.RecipeStepIngredientDatabaseCreationInput]()).Return(expected, nil)
			},
			map[string][]string{
				events.RecipeStepIngredientCreated: {
					keys.RecipeIDKey,
					keys.RecipeStepIDKey,
					keys.RecipeStepIngredientIDKey,
				},
			},
		)

		actual, err := rm.CreateRecipeStepIngredient(ctx, exampleRecipeID, exampleRecipeStepID, fakeInput)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestRecipeManager_ReadRecipeStepIngredient(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		rm := buildRecipeManagerForTest(t)

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()
		expected := fakes.BuildFakeRecipeStepIngredient()

		expectations := setupExpectationsForRecipeManager(
			rm,
			func(db *database.MockDatabase) {
				db.RecipeStepIngredientDataManagerMock.On(testutils.GetMethodName(rm.db.GetRecipeStepIngredient), testutils.ContextMatcher, exampleRecipeID, exampleRecipeStepID, expected.ID).Return(expected, nil)
			},
		)

		actual, err := rm.ReadRecipeStepIngredient(ctx, exampleRecipeID, exampleRecipeStepID, expected.ID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestRecipeManager_UpdateRecipeStepIngredient(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		rm := buildRecipeManagerForTest(t)

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()
		exampleInput := fakes.BuildFakeRecipeStepIngredientUpdateRequestInput()

		expectations := setupExpectationsForRecipeManager(
			rm,
			func(db *database.MockDatabase) {
				db.RecipeStepIngredientDataManagerMock.On(testutils.GetMethodName(rm.db.GetRecipeStepIngredient), testutils.ContextMatcher, exampleRecipeID, exampleRecipeStepID, exampleRecipeStepIngredient.ID).Return(exampleRecipeStepIngredient, nil)
				db.RecipeStepIngredientDataManagerMock.On(testutils.GetMethodName(rm.db.UpdateRecipeStepIngredient), testutils.ContextMatcher, testutils.MatchType[*types.RecipeStepIngredient]()).Return(nil)
			},
			map[string][]string{
				events.RecipeStepIngredientUpdated: {
					keys.RecipeIDKey,
					keys.RecipeStepIDKey,
					keys.RecipeStepIngredientIDKey,
				},
			},
		)

		assert.NoError(t, rm.UpdateRecipeStepIngredient(ctx, exampleRecipeID, exampleRecipeStepID, exampleRecipeStepIngredient.ID, exampleInput))

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestRecipeManager_ArchiveRecipeStepIngredient(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		rm := buildRecipeManagerForTest(t)

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()
		expected := fakes.BuildFakeRecipeStepIngredient()

		expectations := setupExpectationsForRecipeManager(
			rm,
			func(db *database.MockDatabase) {
				db.RecipeStepIngredientDataManagerMock.On(testutils.GetMethodName(rm.db.ArchiveRecipeStepIngredient), testutils.ContextMatcher, exampleRecipeStepID, expected.ID).Return(nil)
			},
			map[string][]string{
				events.RecipeStepIngredientArchived: {
					keys.RecipeIDKey,
					keys.RecipeStepIDKey,
					keys.RecipeStepIngredientIDKey,
				},
			},
		)

		assert.NoError(t, rm.ArchiveRecipeStepIngredient(ctx, exampleRecipeID, exampleRecipeStepID, expected.ID))

		mock.AssertExpectationsForObjects(t, expectations...)
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

		ctx := t.Context()
		rm := buildRecipeManagerForTest(t)

		exampleRecipeID := fakes.BuildFakeID()
		expected := fakes.BuildFakeRecipePrepTask()
		fakeInput := fakes.BuildFakeRecipePrepTaskCreationRequestInput()

		expectations := setupExpectationsForRecipeManager(
			rm,
			func(db *database.MockDatabase) {
				db.RecipePrepTaskDataManagerMock.On(testutils.GetMethodName(rm.db.CreateRecipePrepTask), testutils.ContextMatcher, testutils.MatchType[*types.RecipePrepTaskDatabaseCreationInput]()).Return(expected, nil)
			},
			map[string][]string{
				events.RecipePrepTaskCreated: {
					keys.RecipeIDKey,
					keys.RecipePrepTaskIDKey,
				},
			},
		)

		actual, err := rm.CreateRecipePrepTask(ctx, exampleRecipeID, fakeInput)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestRecipeManager_ReadRecipePrepTask(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		rm := buildRecipeManagerForTest(t)

		exampleRecipeID := fakes.BuildFakeID()
		expected := fakes.BuildFakeRecipePrepTask()

		expectations := setupExpectationsForRecipeManager(
			rm,
			func(db *database.MockDatabase) {
				db.RecipePrepTaskDataManagerMock.On(testutils.GetMethodName(rm.db.GetRecipePrepTask), testutils.ContextMatcher, exampleRecipeID, expected.ID).Return(expected, nil)
			},
		)

		actual, err := rm.ReadRecipePrepTask(ctx, exampleRecipeID, expected.ID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestRecipeManager_UpdateRecipePrepTask(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		rm := buildRecipeManagerForTest(t)

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipePrepTask := fakes.BuildFakeRecipePrepTask()
		exampleInput := fakes.BuildFakeRecipePrepTaskUpdateRequestInput()

		expectations := setupExpectationsForRecipeManager(
			rm,
			func(db *database.MockDatabase) {
				db.RecipePrepTaskDataManagerMock.On(testutils.GetMethodName(rm.db.GetRecipePrepTask), testutils.ContextMatcher, exampleRecipeID, exampleRecipePrepTask.ID).Return(exampleRecipePrepTask, nil)
				db.RecipePrepTaskDataManagerMock.On(testutils.GetMethodName(rm.db.UpdateRecipePrepTask), testutils.ContextMatcher, testutils.MatchType[*types.RecipePrepTask]()).Return(nil)
			},
			map[string][]string{
				events.RecipePrepTaskUpdated: {
					keys.RecipeIDKey,
					keys.RecipePrepTaskIDKey,
				},
			},
		)

		assert.NoError(t, rm.UpdateRecipePrepTask(ctx, exampleRecipeID, exampleRecipePrepTask.ID, exampleInput))

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestRecipeManager_ArchiveRecipePrepTask(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		rm := buildRecipeManagerForTest(t)

		exampleRecipeID := fakes.BuildFakeID()
		expected := fakes.BuildFakeRecipePrepTask()

		expectations := setupExpectationsForRecipeManager(
			rm,
			func(db *database.MockDatabase) {
				db.RecipePrepTaskDataManagerMock.On(testutils.GetMethodName(rm.db.ArchiveRecipePrepTask), testutils.ContextMatcher, exampleRecipeID, expected.ID).Return(nil)
			},
			map[string][]string{
				events.RecipePrepTaskArchived: {
					keys.RecipeIDKey,
					keys.RecipePrepTaskIDKey,
				},
			},
		)

		assert.NoError(t, rm.ArchiveRecipePrepTask(ctx, exampleRecipeID, expected.ID))

		mock.AssertExpectationsForObjects(t, expectations...)
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

		ctx := t.Context()
		rm := buildRecipeManagerForTest(t)

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()
		expected := fakes.BuildFakeRecipeStepCompletionCondition()
		fakeInput := fakes.BuildFakeRecipeStepCompletionConditionForExistingRecipeCreationRequestInput()

		expectations := setupExpectationsForRecipeManager(
			rm,
			func(db *database.MockDatabase) {
				db.RecipeStepCompletionConditionDataManagerMock.On(testutils.GetMethodName(rm.db.CreateRecipeStepCompletionCondition), testutils.ContextMatcher, testutils.MatchType[*types.RecipeStepCompletionConditionDatabaseCreationInput]()).Return(expected, nil)
			},
			map[string][]string{
				events.RecipeStepCompletionConditionCreated: {
					keys.RecipeIDKey,
					keys.RecipeStepIDKey,
					keys.RecipeStepCompletionConditionIDKey,
				},
			},
		)

		actual, err := rm.CreateRecipeStepCompletionCondition(ctx, exampleRecipeID, exampleRecipeStepID, fakeInput)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestRecipeManager_ReadRecipeStepCompletionCondition(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		rm := buildRecipeManagerForTest(t)

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()
		expected := fakes.BuildFakeRecipeStepCompletionCondition()

		expectations := setupExpectationsForRecipeManager(
			rm,
			func(db *database.MockDatabase) {
				db.RecipeStepCompletionConditionDataManagerMock.On(testutils.GetMethodName(rm.db.GetRecipeStepCompletionCondition), testutils.ContextMatcher, exampleRecipeID, exampleRecipeStepID, expected.ID).Return(expected, nil)
			},
		)

		actual, err := rm.ReadRecipeStepCompletionCondition(ctx, exampleRecipeID, exampleRecipeStepID, expected.ID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestRecipeManager_UpdateRecipeStepCompletionCondition(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		rm := buildRecipeManagerForTest(t)

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepCompletionCondition := fakes.BuildFakeRecipeStepCompletionCondition()
		exampleInput := fakes.BuildFakeRecipeStepCompletionConditionUpdateRequestInput()

		expectations := setupExpectationsForRecipeManager(
			rm,
			func(db *database.MockDatabase) {
				db.RecipeStepCompletionConditionDataManagerMock.On(testutils.GetMethodName(rm.db.GetRecipeStepCompletionCondition), testutils.ContextMatcher, exampleRecipeID, exampleRecipeStepID, exampleRecipeStepCompletionCondition.ID).Return(exampleRecipeStepCompletionCondition, nil)
				db.RecipeStepCompletionConditionDataManagerMock.On(testutils.GetMethodName(rm.db.UpdateRecipeStepCompletionCondition), testutils.ContextMatcher, testutils.MatchType[*types.RecipeStepCompletionCondition]()).Return(nil)
			},
			map[string][]string{
				events.RecipeStepCompletionConditionUpdated: {
					keys.RecipeIDKey,
					keys.RecipeStepIDKey,
					keys.RecipeStepCompletionConditionIDKey,
				},
			},
		)

		assert.NoError(t, rm.UpdateRecipeStepCompletionCondition(ctx, exampleRecipeID, exampleRecipeStepID, exampleRecipeStepCompletionCondition.ID, exampleInput))

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestRecipeManager_ArchiveRecipeStepCompletionCondition(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		rm := buildRecipeManagerForTest(t)

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()
		expected := fakes.BuildFakeRecipeStepCompletionCondition()

		expectations := setupExpectationsForRecipeManager(
			rm,
			func(db *database.MockDatabase) {
				db.RecipeStepCompletionConditionDataManagerMock.On(testutils.GetMethodName(rm.db.ArchiveRecipeStepCompletionCondition), testutils.ContextMatcher, exampleRecipeStepID, expected.ID).Return(nil)
			},
			map[string][]string{
				events.RecipeStepCompletionConditionArchived: {
					keys.RecipeIDKey,
					keys.RecipeStepIDKey,
					keys.RecipeStepCompletionConditionIDKey,
				},
			},
		)

		assert.NoError(t, rm.ArchiveRecipeStepCompletionCondition(ctx, exampleRecipeID, exampleRecipeStepID, expected.ID))

		mock.AssertExpectationsForObjects(t, expectations...)
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

		ctx := t.Context()
		rm := buildRecipeManagerForTest(t)

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()
		expected := fakes.BuildFakeRecipeStepVessel()
		fakeInput := fakes.BuildFakeRecipeStepVesselCreationRequestInput()

		expectations := setupExpectationsForRecipeManager(
			rm,
			func(db *database.MockDatabase) {
				db.RecipeStepVesselDataManagerMock.On(testutils.GetMethodName(rm.db.CreateRecipeStepVessel), testutils.ContextMatcher, testutils.MatchType[*types.RecipeStepVesselDatabaseCreationInput]()).Return(expected, nil)
			},
			map[string][]string{
				events.RecipeStepVesselCreated: {
					keys.RecipeIDKey,
					keys.RecipeStepIDKey,
					keys.RecipeStepVesselIDKey,
				},
			},
		)

		actual, err := rm.CreateRecipeStepVessel(ctx, exampleRecipeID, exampleRecipeStepID, fakeInput)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestRecipeManager_ReadRecipeStepVessel(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		rm := buildRecipeManagerForTest(t)

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()
		expected := fakes.BuildFakeRecipeStepVessel()

		expectations := setupExpectationsForRecipeManager(
			rm,
			func(db *database.MockDatabase) {
				db.RecipeStepVesselDataManagerMock.On(testutils.GetMethodName(rm.db.GetRecipeStepVessel), testutils.ContextMatcher, exampleRecipeID, exampleRecipeStepID, expected.ID).Return(expected, nil)
			},
		)

		actual, err := rm.ReadRecipeStepVessel(ctx, exampleRecipeID, exampleRecipeStepID, expected.ID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestRecipeManager_UpdateRecipeStepVessel(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		rm := buildRecipeManagerForTest(t)

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepVessel := fakes.BuildFakeRecipeStepVessel()
		exampleInput := fakes.BuildFakeRecipeStepVesselUpdateRequestInput()

		expectations := setupExpectationsForRecipeManager(
			rm,
			func(db *database.MockDatabase) {
				db.RecipeStepVesselDataManagerMock.On(testutils.GetMethodName(rm.db.GetRecipeStepVessel), testutils.ContextMatcher, exampleRecipeID, exampleRecipeStepID, exampleRecipeStepVessel.ID).Return(exampleRecipeStepVessel, nil)
				db.RecipeStepVesselDataManagerMock.On(testutils.GetMethodName(rm.db.UpdateRecipeStepVessel), testutils.ContextMatcher, testutils.MatchType[*types.RecipeStepVessel]()).Return(nil)
			},
			map[string][]string{
				events.RecipeStepVesselUpdated: {
					keys.RecipeIDKey,
					keys.RecipeStepIDKey,
					keys.RecipeStepVesselIDKey,
				},
			},
		)

		assert.NoError(t, rm.UpdateRecipeStepVessel(ctx, exampleRecipeID, exampleRecipeStepID, exampleRecipeStepVessel.ID, exampleInput))

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestRecipeManager_ArchiveRecipeStepVessel(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		rm := buildRecipeManagerForTest(t)

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()
		expected := fakes.BuildFakeRecipeStepVessel()

		expectations := setupExpectationsForRecipeManager(
			rm,
			func(db *database.MockDatabase) {
				db.RecipeStepVesselDataManagerMock.On(testutils.GetMethodName(rm.db.ArchiveRecipeStepVessel), testutils.ContextMatcher, exampleRecipeStepID, expected.ID).Return(nil)
			},
			map[string][]string{
				events.RecipeStepVesselArchived: {
					keys.RecipeIDKey,
					keys.RecipeStepIDKey,
					keys.RecipeStepVesselIDKey,
				},
			},
		)

		assert.NoError(t, rm.ArchiveRecipeStepVessel(ctx, exampleRecipeID, exampleRecipeStepID, expected.ID))

		mock.AssertExpectationsForObjects(t, expectations...)
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

		ctx := t.Context()
		rm := buildRecipeManagerForTest(t)

		exampleRecipeID := fakes.BuildFakeID()
		expected := fakes.BuildFakeRecipeRating()

		expectations := setupExpectationsForRecipeManager(
			rm,
			func(db *database.MockDatabase) {
				db.RecipeRatingDataManagerMock.On(testutils.GetMethodName(rm.db.GetRecipeRating), testutils.ContextMatcher, exampleRecipeID, expected.ID).Return(expected, nil)
			},
		)

		actual, err := rm.ReadRecipeRating(ctx, exampleRecipeID, expected.ID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestRecipeManager_CreateRecipeRating(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		rm := buildRecipeManagerForTest(t)

		exampleRecipeID := fakes.BuildFakeID()
		expected := fakes.BuildFakeRecipeRating()
		fakeInput := fakes.BuildFakeRecipeRatingCreationRequestInput()

		expectations := setupExpectationsForRecipeManager(
			rm,
			func(db *database.MockDatabase) {
				db.RecipeRatingDataManagerMock.On(testutils.GetMethodName(rm.db.CreateRecipeRating), testutils.ContextMatcher, testutils.MatchType[*types.RecipeRatingDatabaseCreationInput]()).Return(expected, nil)
			},
			map[string][]string{
				events.RecipeRatingCreated: {
					keys.RecipeIDKey,
					keys.RecipeRatingIDKey,
				},
			},
		)

		actual, err := rm.CreateRecipeRating(ctx, exampleRecipeID, fakeInput)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestRecipeManager_UpdateRecipeRating(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		rm := buildRecipeManagerForTest(t)

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeRating := fakes.BuildFakeRecipeRating()
		exampleInput := fakes.BuildFakeRecipeRatingUpdateRequestInput()

		expectations := setupExpectationsForRecipeManager(
			rm,
			func(db *database.MockDatabase) {
				db.RecipeRatingDataManagerMock.On(testutils.GetMethodName(rm.db.GetRecipeRating), testutils.ContextMatcher, exampleRecipeID, exampleRecipeRating.ID).Return(exampleRecipeRating, nil)
				db.RecipeRatingDataManagerMock.On(testutils.GetMethodName(rm.db.UpdateRecipeRating), testutils.ContextMatcher, testutils.MatchType[*types.RecipeRating]()).Return(nil)
			},
			map[string][]string{
				events.RecipeRatingUpdated: {
					keys.RecipeIDKey,
					keys.RecipeRatingIDKey,
				},
			},
		)

		assert.NoError(t, rm.UpdateRecipeRating(ctx, exampleRecipeID, exampleRecipeRating.ID, exampleInput))

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestRecipeManager_ArchiveRecipeRating(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		rm := buildRecipeManagerForTest(t)

		exampleRecipeID := fakes.BuildFakeID()
		expected := fakes.BuildFakeRecipeRating()

		expectations := setupExpectationsForRecipeManager(
			rm,
			func(db *database.MockDatabase) {
				db.RecipeRatingDataManagerMock.On(testutils.GetMethodName(rm.db.ArchiveRecipeRating), testutils.ContextMatcher, exampleRecipeID, expected.ID).Return(nil)
			},
			map[string][]string{
				events.RecipeRatingArchived: {
					keys.RecipeIDKey,
					keys.RecipeRatingIDKey,
				},
			},
		)

		assert.NoError(t, rm.ArchiveRecipeRating(ctx, exampleRecipeID, expected.ID))

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}
