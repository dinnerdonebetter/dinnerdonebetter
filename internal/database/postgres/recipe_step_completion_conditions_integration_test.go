package postgres

import (
	"context"
	"database/sql"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/identifiers"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createRecipeStepCompletionConditionForTest(t *testing.T, ctx context.Context, recipeID string, exampleRecipeStepCompletionCondition *types.RecipeStepCompletionCondition, dbc *Querier) *types.RecipeStepCompletionCondition {
	t.Helper()

	// create
	if exampleRecipeStepCompletionCondition == nil {
		t.Fatal("exampleRecipeStepCompletionCondition must not be nil")
	}
	dbInput := converters.ConvertRecipeStepCompletionConditionToRecipeStepCompletionConditionDatabaseCreationInput(exampleRecipeStepCompletionCondition)

	recipe, getRecipeErr := dbc.GetRecipe(ctx, recipeID)
	require.NoError(t, getRecipeErr)
	require.NotNil(t, recipe)

	recipeStep, getRecipeStepErr := dbc.GetRecipeStep(ctx, recipeID, exampleRecipeStepCompletionCondition.BelongsToRecipeStep)
	require.NoError(t, getRecipeStepErr)
	require.NotNil(t, recipeStep)

	ingredientState, getIngredientStateErr := dbc.GetValidIngredientState(ctx, dbInput.IngredientStateID)
	require.NoError(t, getIngredientStateErr)
	require.NotNil(t, ingredientState)

	created, err := dbc.CreateRecipeStepCompletionCondition(ctx, dbInput)
	assert.NoError(t, err)
	require.NotNil(t, created)

	exampleRecipeStepCompletionCondition.CreatedAt = created.CreatedAt
	exampleRecipeStepCompletionCondition.IngredientState = created.IngredientState
	exampleRecipeStepCompletionCondition.Ingredients = created.Ingredients
	assert.Equal(t, exampleRecipeStepCompletionCondition, created)

	recipeStepCompletionCondition, err := dbc.GetRecipeStepCompletionCondition(ctx, recipe.ID, recipeStep.ID, created.ID)
	require.NoError(t, err)
	require.NotNil(t, recipeStepCompletionCondition)

	exampleRecipeStepCompletionCondition.CreatedAt = recipeStepCompletionCondition.CreatedAt
	exampleRecipeStepCompletionCondition.IngredientState = recipeStepCompletionCondition.IngredientState
	exampleRecipeStepCompletionCondition.Ingredients = recipeStepCompletionCondition.Ingredients

	require.Equal(t, exampleRecipeStepCompletionCondition, recipeStepCompletionCondition)

	return recipeStepCompletionCondition
}

func TestQuerier_Integration_RecipeStepCompletionConditions(t *testing.T) {
	// if !runningContainerTests {
	t.SkipNow()
	// }

	ctx := context.Background()
	dbc, container := buildDatabaseClientForTest(t, ctx)

	databaseURI, err := container.ConnectionString(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, databaseURI)

	defer func(t *testing.T) {
		t.Helper()
		assert.NoError(t, container.Terminate(ctx))
	}(t)

	user := createUserForTest(t, ctx, nil, dbc)
	householdID, err := dbc.GetDefaultHouseholdIDForUser(ctx, user.ID)
	require.NoError(t, err)
	require.NotEmpty(t, householdID)

	exampleRecipe := buildRecipeForTestCreation(t, ctx, user.ID, dbc)
	createdRecipe := createRecipeForTest(t, ctx, exampleRecipe, dbc, true)
	exampleRecipeStep := createdRecipe.Steps[0]

	ingredientState := createValidIngredientStateForTest(t, ctx, nil, dbc)
	exampleRecipeStepCompletionCondition := fakes.BuildFakeRecipeStepCompletionCondition()
	exampleRecipeStepCompletionCondition.IngredientState = *ingredientState
	exampleRecipeStepCompletionCondition.BelongsToRecipeStep = exampleRecipeStep.ID
	exampleRecipeStepCompletionCondition.Ingredients = []*types.RecipeStepCompletionConditionIngredient{
		{
			ID:                                     identifiers.New(),
			BelongsToRecipeStepCompletionCondition: exampleRecipeStepCompletionCondition.ID,
			RecipeStepIngredient:                   exampleRecipeStep.Ingredients[0].ID,
		},
	}

	// create
	createdRecipeStepCompletionConditions := []*types.RecipeStepCompletionCondition{}
	createdRecipeStepCompletionConditions = append(createdRecipeStepCompletionConditions, createRecipeStepCompletionConditionForTest(t, ctx, exampleRecipe.ID, exampleRecipeStepCompletionCondition, dbc))

	// fetch as list
	recipeStepCompletionConditions, err := dbc.GetRecipeStepCompletionConditions(ctx, exampleRecipe.ID, createdRecipeStepCompletionConditions[0].BelongsToRecipeStep, nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, recipeStepCompletionConditions.Data)
	assert.Equal(t, len(createdRecipeStepCompletionConditions), len(recipeStepCompletionConditions.Data))

	// delete
	for _, recipeStepCompletionCondition := range createdRecipeStepCompletionConditions {
		assert.NoError(t, dbc.ArchiveRecipeStepCompletionCondition(ctx, exampleRecipeStep.ID, recipeStepCompletionCondition.ID))

		var exists bool
		exists, err = dbc.RecipeStepCompletionConditionExists(ctx, exampleRecipe.ID, recipeStepCompletionCondition.BelongsToRecipeStep, recipeStepCompletionCondition.ID)
		assert.NoError(t, err)
		assert.False(t, exists)

		var y *types.RecipeStepCompletionCondition
		y, err = dbc.GetRecipeStepCompletionCondition(ctx, exampleRecipe.ID, recipeStepCompletionCondition.BelongsToRecipeStep, recipeStepCompletionCondition.ID)
		assert.Nil(t, y)
		assert.Error(t, err)
		assert.ErrorIs(t, err, sql.ErrNoRows)
	}
}
