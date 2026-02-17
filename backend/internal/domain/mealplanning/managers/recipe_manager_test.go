package managers

import (
	"errors"
	"testing"

	types "github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning/fakes"
	mealplanningmock "github.com/dinnerdonebetter/backend/internal/domain/mealplanning/mocks"
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning/recipeanalysis"
	"github.com/dinnerdonebetter/backend/internal/platform/database/filtering"
	msgconfig "github.com/dinnerdonebetter/backend/internal/platform/messagequeue/config"
	mockpublishers "github.com/dinnerdonebetter/backend/internal/platform/messagequeue/mock"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/metrics"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/platform/reflection"
	textsearchcfg "github.com/dinnerdonebetter/backend/internal/platform/search/text/config"
	"github.com/dinnerdonebetter/backend/internal/platform/testutils"

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
	mpp.On(reflection.GetMethodName(mpp.ProvidePublisher), queueCfg.DataChangesTopicName).Return(&mockpublishers.Publisher{}, nil)

	m, err := NewRecipeManager(
		t.Context(),
		logging.NewNoopLogger(),
		tracing.NewNoopTracerProvider(),
		&mealplanningmock.Repository{},
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
	dbSetupFunc func(db *mealplanningmock.Repository),
	eventTypeMaps ...map[string][]string,
) []any {
	db := &mealplanningmock.Repository{}
	if dbSetupFunc != nil {
		dbSetupFunc(db)
	}
	manager.db = db

	mp := &mockpublishers.Publisher{}
	for _, eventTypeMap := range eventTypeMaps {
		for eventType, payload := range eventTypeMap {
			mp.On(reflection.GetMethodName(mp.PublishAsync), testutils.ContextMatcher, eventMatches(eventType, payload)).Return()
		}
	}
	manager.dataChangesPublisher = mp

	return []any{db, mp}
}

func setupExpectationsForRecipeManagerWithAnalyzer(
	manager *recipeManager,
	dbSetupFunc func(db *mealplanningmock.Repository),
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
		status := types.RecipeStatusSubmitted

		expectations := setupExpectationsForRecipeManager(
			rm,
			func(db *mealplanningmock.Repository) {
				db.On(reflection.GetMethodName(rm.db.GetRecipes), testutils.ContextMatcher, status, testutils.QueryFilterMatcher).Return(expected, nil)
			},
		)

		actual, err := rm.ListRecipes(ctx, status, nil)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestRecipeManager_CreateRecipe(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		rm := buildRecipeManagerForTest(t)

		fakeCreatorID := fakes.BuildFakeID()
		expected := fakes.BuildFakeRecipe()
		fakeInput := fakes.BuildFakeRecipeCreationRequestInput()

		analyzer := &recipeanalysis.MockRecipeAnalyzer{}
		analyzer.On(reflection.GetMethodName(analyzer.ValidateRecipeCreationRequestInputIsDAG), testutils.ContextMatcher, testutils.MatchType[*types.RecipeCreationRequestInput]()).Return(nil)
		rm.recipeAnalyzer = analyzer

		expectations := setupExpectationsForRecipeManager(
			rm,
			func(db *mealplanningmock.Repository) {
				db.On(reflection.GetMethodName(rm.db.CreateRecipe), testutils.ContextMatcher, testutils.MatchType[*types.RecipeDatabaseCreationInput]()).Return(expected, nil)
				db.On(reflection.GetMethodName(rm.db.GetRecipe), testutils.ContextMatcher, expected.ID).Return(expected, nil)
			},
			map[string][]string{
				types.RecipeCreatedServiceEventType: {
					keys.RecipeIDKey,
				},
			},
		)

		actual, err := rm.CreateRecipe(ctx, fakeCreatorID, fakeInput)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mock.AssertExpectationsForObjects(t, append(expectations, analyzer)...)
	})

	T.Run("with DAG error", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		rm := buildRecipeManagerForTest(t)

		fakeCreatorID := fakes.BuildFakeID()
		fakeInput := fakes.BuildFakeRecipeCreationRequestInput()

		analyzer := &recipeanalysis.MockRecipeAnalyzer{}
		analyzer.On(reflection.GetMethodName(analyzer.ValidateRecipeCreationRequestInputIsDAG), testutils.ContextMatcher, testutils.MatchType[*types.RecipeCreationRequestInput]()).Return(errors.New("blah"))
		rm.recipeAnalyzer = analyzer

		actual, err := rm.CreateRecipe(ctx, fakeCreatorID, fakeInput)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, analyzer)
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
			func(db *mealplanningmock.Repository) {
				db.On(reflection.GetMethodName(rm.db.GetRecipe), testutils.ContextMatcher, expected.ID).Return(expected, nil)
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
			func(db *mealplanningmock.Repository) {
				db.On(reflection.GetMethodName(rm.db.SearchForRecipes), testutils.ContextMatcher, exampleQuery, testutils.QueryFilterMatcher).Return(expected, nil)
			},
		)

		actual, err := rm.SearchRecipes(ctx, exampleQuery, false, nil)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

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
			func(db *mealplanningmock.Repository) {
				db.On(reflection.GetMethodName(rm.db.GetRecipe), testutils.ContextMatcher, exampleRecipe.ID).Return(exampleRecipe, nil)
				db.On(reflection.GetMethodName(rm.db.UpdateRecipe), testutils.ContextMatcher, testutils.MatchType[*types.Recipe]()).Return(nil)
			},
			map[string][]string{
				types.RecipeUpdatedServiceEventType: {
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
			func(db *mealplanningmock.Repository) {
				db.On(reflection.GetMethodName(rm.db.ArchiveRecipe), testutils.ContextMatcher, expected.ID, exampleOwnerID).Return(nil)
			},
			map[string][]string{
				types.RecipeArchivedServiceEventType: {
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
			func(db *mealplanningmock.Repository) {
				db.On(reflection.GetMethodName(rm.db.GetRecipe), testutils.ContextMatcher, exampleRecipe.ID).Return(exampleRecipe, nil)
			},
			func(analyzer *recipeanalysis.MockRecipeAnalyzer) {
				analyzer.On(reflection.GetMethodName(analyzer.GenerateMealPlanTasksForRecipe), testutils.ContextMatcher, "", testutils.MatchType[*types.Recipe]()).Return(expectedResults, nil)
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
			func(db *mealplanningmock.Repository) {
				db.On(reflection.GetMethodName(rm.db.GetRecipe), testutils.ContextMatcher, exampleRecipe.ID).Return(exampleRecipe, nil)
			},
			func(analyzer *recipeanalysis.MockRecipeAnalyzer) {
				analyzer.On(reflection.GetMethodName(analyzer.RenderMermaidDiagramForRecipe), testutils.ContextMatcher, testutils.MatchType[*types.Recipe]()).Return(expectedResult, nil)
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
			func(db *mealplanningmock.Repository) {
				db.On(reflection.GetMethodName(rm.db.GetRecipe), testutils.ContextMatcher, expected.ID).Return(expected, nil)
				db.On(reflection.GetMethodName(rm.db.CreateRecipe), testutils.ContextMatcher, testutils.MatchType[*types.RecipeDatabaseCreationInput]()).Return(cloned, nil)
			},
			map[string][]string{
				types.RecipeClonedServiceEventType: {
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

func TestRecipeManager_ListRecipeLists(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		rm := buildRecipeManagerForTest(t)

		recipeList := &types.RecipeList{
			ID:            fakes.BuildFakeID(),
			Name:          t.Name(),
			Description:   t.Name(),
			BelongsToUser: fakes.BuildFakeID(),
		}
		expected := &filtering.QueryFilteredResult[types.RecipeList]{Data: []*types.RecipeList{recipeList}}

		expectations := setupExpectationsForRecipeManager(
			rm,
			func(db *mealplanningmock.Repository) {
				db.On(reflection.GetMethodName(rm.db.GetRecipeLists), testutils.ContextMatcher, testutils.QueryFilterMatcher).Return(expected, nil)
			},
		)

		actual, err := rm.ListRecipeLists(ctx, nil)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestRecipeManager_CreateRecipeList(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		rm := buildRecipeManagerForTest(t)

		userID := fakes.BuildFakeID()
		input := &types.RecipeListCreationRequestInput{
			Name:        t.Name(),
			Description: t.Name(),
		}
		expected := &types.RecipeList{ID: fakes.BuildFakeID(), Name: input.Name, Description: input.Description, BelongsToUser: userID}

		expectations := setupExpectationsForRecipeManager(
			rm,
			func(db *mealplanningmock.Repository) {
				db.On(reflection.GetMethodName(rm.db.CreateRecipeList), testutils.ContextMatcher, testutils.MatchType[*types.RecipeListDatabaseCreationInput]()).Return(expected, nil)
			},
		)

		actual, err := rm.CreateRecipeList(ctx, userID, input)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestRecipeManager_ArchiveRecipeList(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		rm := buildRecipeManagerForTest(t)

		userID := fakes.BuildFakeID()
		listID := fakes.BuildFakeID()

		expectations := setupExpectationsForRecipeManager(
			rm,
			func(db *mealplanningmock.Repository) {
				db.On(reflection.GetMethodName(rm.db.ArchiveRecipeList), testutils.ContextMatcher, listID, userID).Return(nil)
			},
		)

		assert.NoError(t, rm.ArchiveRecipeList(ctx, listID, userID))

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestRecipeManager_UpdateRecipeList(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		rm := buildRecipeManagerForTest(t)

		listID := fakes.BuildFakeID()
		userID := fakes.BuildFakeID()
		name := t.Name()
		desc := "desc"
		input := &types.RecipeListUpdateRequestInput{
			Name:        &name,
			Description: &desc,
		}

		expectations := setupExpectationsForRecipeManager(
			rm,
			func(db *mealplanningmock.Repository) {
				db.On(reflection.GetMethodName(rm.db.UpdateRecipeList), testutils.ContextMatcher, testutils.MatchType[*types.RecipeList]()).Return(nil)
			},
		)

		assert.NoError(t, rm.UpdateRecipeList(ctx, listID, userID, input))

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestRecipeManager_UpdateRecipeListItem(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		rm := buildRecipeManagerForTest(t)

		itemID := fakes.BuildFakeID()
		listID := fakes.BuildFakeID()
		recipeID := fakes.BuildFakeID()
		notes := new(t.Name())
		input := &types.RecipeListItemUpdateRequestInput{
			Notes: notes,
		}

		expectations := setupExpectationsForRecipeManager(
			rm,
			func(db *mealplanningmock.Repository) {
				db.On(reflection.GetMethodName(rm.db.UpdateRecipeListItem), testutils.ContextMatcher, testutils.MatchType[*types.RecipeListItem]()).Return(nil)
			},
		)

		assert.NoError(t, rm.UpdateRecipeListItem(ctx, itemID, listID, recipeID, input))

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestRecipeManager_AddRecipeToRecipeList(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		rm := buildRecipeManagerForTest(t)

		listID := fakes.BuildFakeID()
		recipeID := fakes.BuildFakeID()
		expected := &types.RecipeListItem{
			ID:                  fakes.BuildFakeID(),
			BelongsToRecipeList: listID,
			Notes:               t.Name(),
			Recipe:              types.Recipe{ID: recipeID},
		}

		expectations := setupExpectationsForRecipeManager(
			rm,
			func(db *mealplanningmock.Repository) {
				db.On(reflection.GetMethodName(rm.db.CreateRecipeListItem), testutils.ContextMatcher, testutils.MatchType[*types.RecipeListItemDatabaseCreationInput]()).Return(expected, nil)
			},
		)

		actual, err := rm.AddRecipeToRecipeList(ctx, listID, recipeID, expected.Notes)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestRecipeManager_RemoveRecipeFromRecipeList(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		rm := buildRecipeManagerForTest(t)

		listID := fakes.BuildFakeID()
		itemID := fakes.BuildFakeID()

		expectations := setupExpectationsForRecipeManager(
			rm,
			func(db *mealplanningmock.Repository) {
				db.On(reflection.GetMethodName(rm.db.ArchiveRecipeListItem), testutils.ContextMatcher, itemID, listID).Return(nil)
			},
		)

		assert.NoError(t, rm.RemoveRecipeFromRecipeList(ctx, listID, itemID))

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestRecipeManager_ListRecipeListItems(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		rm := buildRecipeManagerForTest(t)

		listID := fakes.BuildFakeID()
		expectedItem := &types.RecipeListItem{
			ID:                  fakes.BuildFakeID(),
			BelongsToRecipeList: listID,
			Notes:               t.Name(),
			Recipe:              types.Recipe{ID: fakes.BuildFakeID()},
		}
		expected := &filtering.QueryFilteredResult[types.RecipeListItem]{Data: []*types.RecipeListItem{expectedItem}}

		expectations := setupExpectationsForRecipeManager(
			rm,
			func(db *mealplanningmock.Repository) {
				db.On(reflection.GetMethodName(rm.db.GetRecipeListItems), testutils.ContextMatcher, listID, testutils.QueryFilterMatcher).Return(expected, nil)
			},
		)

		actual, err := rm.ListRecipeListItems(ctx, listID, nil)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

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
			func(db *mealplanningmock.Repository) {
				db.On(reflection.GetMethodName(rm.db.GetRecipeSteps), testutils.ContextMatcher, exampleRecipeID, testutils.QueryFilterMatcher).Return(expected, nil)
			},
		)

		actual, err := rm.ListRecipeSteps(ctx, exampleRecipeID, nil)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

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
			func(db *mealplanningmock.Repository) {
				db.On(reflection.GetMethodName(rm.db.CreateRecipeStep), testutils.ContextMatcher, testutils.MatchType[*types.RecipeStepDatabaseCreationInput]()).Return(expected, nil)
			},
			map[string][]string{
				types.RecipeStepCreatedServiceEventType: {
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
			func(db *mealplanningmock.Repository) {
				db.On(reflection.GetMethodName(rm.db.GetRecipeStep), testutils.ContextMatcher, exampleRecipeID, expected.ID).Return(expected, nil)
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
			func(db *mealplanningmock.Repository) {
				db.On(reflection.GetMethodName(rm.db.GetRecipeStep), testutils.ContextMatcher, exampleRecipeID, exampleRecipeStep.ID).Return(exampleRecipeStep, nil)
				db.On(reflection.GetMethodName(rm.db.UpdateRecipeStep), testutils.ContextMatcher, testutils.MatchType[*types.RecipeStep]()).Return(nil)
			},
			map[string][]string{
				types.RecipeStepUpdatedServiceEventType: {
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
			func(db *mealplanningmock.Repository) {
				db.On(reflection.GetMethodName(rm.db.ArchiveRecipeStep), testutils.ContextMatcher, exampleRecipeID, expected.ID).Return(nil)
			},
			map[string][]string{
				types.RecipeStepArchivedServiceEventType: {
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
			func(db *mealplanningmock.Repository) {
				db.On(reflection.GetMethodName(rm.db.GetRecipeStepProducts), testutils.ContextMatcher, exampleRecipeID, exampleRecipeStepID, testutils.QueryFilterMatcher).Return(expected, nil)
			},
		)

		actual, err := rm.ListRecipeStepProducts(ctx, exampleRecipeID, exampleRecipeStepID, nil)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

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
			func(db *mealplanningmock.Repository) {
				db.On(reflection.GetMethodName(rm.db.CreateRecipeStepProduct), testutils.ContextMatcher, testutils.MatchType[*types.RecipeStepProductDatabaseCreationInput]()).Return(expected, nil)
			},
			map[string][]string{
				types.RecipeStepProductCreatedServiceEventType: {
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
			func(db *mealplanningmock.Repository) {
				db.On(reflection.GetMethodName(rm.db.GetRecipeStepProduct), testutils.ContextMatcher, exampleRecipeID, exampleRecipeStepID, expected.ID).Return(expected, nil)
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
			func(db *mealplanningmock.Repository) {
				db.On(reflection.GetMethodName(rm.db.GetRecipeStepProduct), testutils.ContextMatcher, exampleRecipeID, exampleRecipeStepID, exampleRecipeStepProduct.ID).Return(exampleRecipeStepProduct, nil)
				db.On(reflection.GetMethodName(rm.db.UpdateRecipeStepProduct), testutils.ContextMatcher, testutils.MatchType[*types.RecipeStepProduct]()).Return(nil)
			},
			map[string][]string{
				types.RecipeStepProductUpdatedServiceEventType: {
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
			func(db *mealplanningmock.Repository) {
				db.On(reflection.GetMethodName(rm.db.ArchiveRecipeStepProduct), testutils.ContextMatcher, exampleRecipeStepID, expected.ID).Return(nil)
			},
			map[string][]string{
				types.RecipeStepProductArchivedServiceEventType: {
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
			func(db *mealplanningmock.Repository) {
				db.On(reflection.GetMethodName(rm.db.GetRecipeStepInstruments), testutils.ContextMatcher, exampleRecipeID, exampleRecipeStepID, testutils.QueryFilterMatcher).Return(expected, nil)
			},
		)

		actual, err := rm.ListRecipeStepInstruments(ctx, exampleRecipeID, exampleRecipeStepID, nil)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

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
		fakeInput.Index = new(uint16(0))

		// Create a fake ValidPreparationInstrument for the bridge table lookup
		fakeValidPreparationInstrument := fakes.BuildFakeValidPreparationInstrument()

		expectations := setupExpectationsForRecipeManager(
			rm,
			func(db *mealplanningmock.Repository) {
				db.On(reflection.GetMethodName(rm.db.GetValidPreparationInstrument), testutils.ContextMatcher, *fakeInput.ValidPreparationInstrumentID).Return(fakeValidPreparationInstrument, nil)
				db.On(reflection.GetMethodName(rm.db.CreateRecipeStepInstrument), testutils.ContextMatcher, testutils.MatchType[*types.RecipeStepInstrumentDatabaseCreationInput]()).Return(expected, nil)
			},
			map[string][]string{
				types.RecipeStepInstrumentCreatedServiceEventType: {
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
			func(db *mealplanningmock.Repository) {
				db.On(reflection.GetMethodName(rm.db.GetRecipeStepInstrument), testutils.ContextMatcher, exampleRecipeID, exampleRecipeStepID, expected.ID).Return(expected, nil)
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
			func(db *mealplanningmock.Repository) {
				db.On(reflection.GetMethodName(rm.db.GetRecipeStepInstrument), testutils.ContextMatcher, exampleRecipeID, exampleRecipeStepID, exampleRecipeStepInstrument.ID).Return(exampleRecipeStepInstrument, nil)
				db.On(reflection.GetMethodName(rm.db.UpdateRecipeStepInstrument), testutils.ContextMatcher, testutils.MatchType[*types.RecipeStepInstrument]()).Return(nil)
			},
			map[string][]string{
				types.RecipeStepInstrumentUpdatedServiceEventType: {
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
			func(db *mealplanningmock.Repository) {
				db.On(reflection.GetMethodName(rm.db.ArchiveRecipeStepInstrument), testutils.ContextMatcher, exampleRecipeStepID, expected.ID).Return(nil)
			},
			map[string][]string{
				types.RecipeStepInstrumentArchivedServiceEventType: {
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
			func(db *mealplanningmock.Repository) {
				db.On(reflection.GetMethodName(rm.db.GetRecipeStepIngredients), testutils.ContextMatcher, exampleRecipeID, exampleRecipeStepID, testutils.QueryFilterMatcher).Return(expected, nil)
			},
		)

		actual, err := rm.ListRecipeStepIngredients(ctx, exampleRecipeID, exampleRecipeStepID, nil)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

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
		fakeInput.Index = new(uint16(0))

		// Create fake bridge table entries for the lookups
		fakeValidIngredientPreparation := fakes.BuildFakeValidIngredientPreparation()
		fakeValidIngredientMeasurementUnit := fakes.BuildFakeValidIngredientMeasurementUnit()

		expectations := setupExpectationsForRecipeManager(
			rm,
			func(db *mealplanningmock.Repository) {
				db.On(reflection.GetMethodName(rm.db.GetValidIngredientPreparation), testutils.ContextMatcher, *fakeInput.ValidIngredientPreparationID).Return(fakeValidIngredientPreparation, nil)
				db.On(reflection.GetMethodName(rm.db.GetValidIngredientMeasurementUnit), testutils.ContextMatcher, *fakeInput.ValidIngredientMeasurementUnitID).Return(fakeValidIngredientMeasurementUnit, nil)
				db.On(reflection.GetMethodName(rm.db.CreateRecipeStepIngredient), testutils.ContextMatcher, testutils.MatchType[*types.RecipeStepIngredientDatabaseCreationInput]()).Return(expected, nil)
			},
			map[string][]string{
				types.RecipeStepIngredientCreatedServiceEventType: {
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
			func(db *mealplanningmock.Repository) {
				db.On(reflection.GetMethodName(rm.db.GetRecipeStepIngredient), testutils.ContextMatcher, exampleRecipeID, exampleRecipeStepID, expected.ID).Return(expected, nil)
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
			func(db *mealplanningmock.Repository) {
				db.On(reflection.GetMethodName(rm.db.GetRecipeStepIngredient), testutils.ContextMatcher, exampleRecipeID, exampleRecipeStepID, exampleRecipeStepIngredient.ID).Return(exampleRecipeStepIngredient, nil)
				db.On(reflection.GetMethodName(rm.db.UpdateRecipeStepIngredient), testutils.ContextMatcher, testutils.MatchType[*types.RecipeStepIngredient]()).Return(nil)
			},
			map[string][]string{
				types.RecipeStepIngredientUpdatedServiceEventType: {
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
			func(db *mealplanningmock.Repository) {
				db.On(reflection.GetMethodName(rm.db.ArchiveRecipeStepIngredient), testutils.ContextMatcher, exampleRecipeStepID, expected.ID).Return(nil)
			},
			map[string][]string{
				types.RecipeStepIngredientArchivedServiceEventType: {
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
			func(db *mealplanningmock.Repository) {
				db.On(reflection.GetMethodName(rm.db.GetRecipePrepTasks), testutils.ContextMatcher, exampleRecipeID, testutils.QueryFilterMatcher).Return(expected, nil)
			},
		)

		actual, err := rm.ListRecipePrepTask(ctx, exampleRecipeID, nil)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

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
			func(db *mealplanningmock.Repository) {
				db.On(reflection.GetMethodName(rm.db.CreateRecipePrepTask), testutils.ContextMatcher, testutils.MatchType[*types.RecipePrepTaskDatabaseCreationInput]()).Return(expected, nil)
			},
			map[string][]string{
				types.RecipePrepTaskCreatedServiceEventType: {
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
			func(db *mealplanningmock.Repository) {
				db.On(reflection.GetMethodName(rm.db.GetRecipePrepTask), testutils.ContextMatcher, exampleRecipeID, expected.ID).Return(expected, nil)
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
			func(db *mealplanningmock.Repository) {
				db.On(reflection.GetMethodName(rm.db.GetRecipePrepTask), testutils.ContextMatcher, exampleRecipeID, exampleRecipePrepTask.ID).Return(exampleRecipePrepTask, nil)
				db.On(reflection.GetMethodName(rm.db.UpdateRecipePrepTask), testutils.ContextMatcher, testutils.MatchType[*types.RecipePrepTask]()).Return(nil)
			},
			map[string][]string{
				types.RecipePrepTaskUpdatedServiceEventType: {
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
			func(db *mealplanningmock.Repository) {
				db.On(reflection.GetMethodName(rm.db.ArchiveRecipePrepTask), testutils.ContextMatcher, exampleRecipeID, expected.ID).Return(nil)
			},
			map[string][]string{
				types.RecipePrepTaskArchivedServiceEventType: {
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
			func(db *mealplanningmock.Repository) {
				db.On(reflection.GetMethodName(rm.db.GetRecipeStepCompletionConditions), testutils.ContextMatcher, exampleRecipeID, exampleRecipeStepID, testutils.QueryFilterMatcher).Return(expected, nil)
			},
		)

		actual, err := rm.ListRecipeStepCompletionConditions(ctx, exampleRecipeID, exampleRecipeStepID, nil)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

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
			func(db *mealplanningmock.Repository) {
				db.On(reflection.GetMethodName(rm.db.CreateRecipeStepCompletionCondition), testutils.ContextMatcher, testutils.MatchType[*types.RecipeStepCompletionConditionDatabaseCreationInput]()).Return(expected, nil)
			},
			map[string][]string{
				types.RecipeStepCompletionConditionCreatedServiceEventType: {
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
			func(db *mealplanningmock.Repository) {
				db.On(reflection.GetMethodName(rm.db.GetRecipeStepCompletionCondition), testutils.ContextMatcher, exampleRecipeID, exampleRecipeStepID, expected.ID).Return(expected, nil)
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
			func(db *mealplanningmock.Repository) {
				db.On(reflection.GetMethodName(rm.db.GetRecipeStepCompletionCondition), testutils.ContextMatcher, exampleRecipeID, exampleRecipeStepID, exampleRecipeStepCompletionCondition.ID).Return(exampleRecipeStepCompletionCondition, nil)
				db.On(reflection.GetMethodName(rm.db.UpdateRecipeStepCompletionCondition), testutils.ContextMatcher, testutils.MatchType[*types.RecipeStepCompletionCondition]()).Return(nil)
			},
			map[string][]string{
				types.RecipeStepCompletionConditionUpdatedServiceEventType: {
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
			func(db *mealplanningmock.Repository) {
				db.On(reflection.GetMethodName(rm.db.ArchiveRecipeStepCompletionCondition), testutils.ContextMatcher, exampleRecipeStepID, expected.ID).Return(nil)
			},
			map[string][]string{
				types.RecipeStepCompletionConditionArchivedServiceEventType: {
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
			func(db *mealplanningmock.Repository) {
				db.On(reflection.GetMethodName(rm.db.GetRecipeStepVessels), testutils.ContextMatcher, exampleRecipeID, exampleRecipeStepID, testutils.QueryFilterMatcher).Return(expected, nil)
			},
		)

		actual, err := rm.ListRecipeStepVessels(ctx, exampleRecipeID, exampleRecipeStepID, nil)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

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
		fakeInput.Index = new(uint16(0))

		// Create a fake ValidPreparationVessel for the bridge table lookup
		fakeValidPreparationVessel := fakes.BuildFakeValidPreparationVessel()

		expectations := setupExpectationsForRecipeManager(
			rm,
			func(db *mealplanningmock.Repository) {
				db.On(reflection.GetMethodName(rm.db.GetValidPreparationVessel), testutils.ContextMatcher, *fakeInput.ValidPreparationVesselID).Return(fakeValidPreparationVessel, nil)
				db.On(reflection.GetMethodName(rm.db.CreateRecipeStepVessel), testutils.ContextMatcher, testutils.MatchType[*types.RecipeStepVesselDatabaseCreationInput]()).Return(expected, nil)
			},
			map[string][]string{
				types.RecipeStepVesselCreatedServiceEventType: {
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
			func(db *mealplanningmock.Repository) {
				db.On(reflection.GetMethodName(rm.db.GetRecipeStepVessel), testutils.ContextMatcher, exampleRecipeID, exampleRecipeStepID, expected.ID).Return(expected, nil)
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
			func(db *mealplanningmock.Repository) {
				db.On(reflection.GetMethodName(rm.db.GetRecipeStepVessel), testutils.ContextMatcher, exampleRecipeID, exampleRecipeStepID, exampleRecipeStepVessel.ID).Return(exampleRecipeStepVessel, nil)
				db.On(reflection.GetMethodName(rm.db.UpdateRecipeStepVessel), testutils.ContextMatcher, testutils.MatchType[*types.RecipeStepVessel]()).Return(nil)
			},
			map[string][]string{
				types.RecipeStepVesselUpdatedServiceEventType: {
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
			func(db *mealplanningmock.Repository) {
				db.On(reflection.GetMethodName(rm.db.ArchiveRecipeStepVessel), testutils.ContextMatcher, exampleRecipeStepID, expected.ID).Return(nil)
			},
			map[string][]string{
				types.RecipeStepVesselArchivedServiceEventType: {
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
			func(db *mealplanningmock.Repository) {
				db.On(reflection.GetMethodName(rm.db.GetRecipeRatingsForRecipe), testutils.ContextMatcher, exampleRecipeID, testutils.QueryFilterMatcher).Return(expected, nil)
			},
		)

		actual, err := rm.ListRecipeRatings(ctx, exampleRecipeID, nil)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

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
			func(db *mealplanningmock.Repository) {
				db.On(reflection.GetMethodName(rm.db.GetRecipeRating), testutils.ContextMatcher, exampleRecipeID, expected.ID).Return(expected, nil)
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
			func(db *mealplanningmock.Repository) {
				db.On(reflection.GetMethodName(rm.db.CreateRecipeRating), testutils.ContextMatcher, testutils.MatchType[*types.RecipeRatingDatabaseCreationInput]()).Return(expected, nil)
			},
			map[string][]string{
				types.RecipeRatingCreatedServiceEventType: {
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
			func(db *mealplanningmock.Repository) {
				db.On(reflection.GetMethodName(rm.db.GetRecipeRating), testutils.ContextMatcher, exampleRecipeID, exampleRecipeRating.ID).Return(exampleRecipeRating, nil)
				db.On(reflection.GetMethodName(rm.db.UpdateRecipeRating), testutils.ContextMatcher, testutils.MatchType[*types.RecipeRating]()).Return(nil)
			},
			map[string][]string{
				types.RecipeRatingUpdatedServiceEventType: {
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
			func(db *mealplanningmock.Repository) {
				db.On(reflection.GetMethodName(rm.db.ArchiveRecipeRating), testutils.ContextMatcher, exampleRecipeID, expected.ID).Return(nil)
			},
			map[string][]string{
				types.RecipeRatingArchivedServiceEventType: {
					keys.RecipeIDKey,
					keys.RecipeRatingIDKey,
				},
			},
		)

		assert.NoError(t, rm.ArchiveRecipeRating(ctx, exampleRecipeID, expected.ID))

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}
