package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/identifiers"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func buildRecipeForTestCreation(t *testing.T, ctx context.Context, userID string, dbc *Querier) *types.Recipe {
	t.Helper()

	if userID == "" {
		user := createUserForTest(t, ctx, nil, dbc)
		userID = user.ID
	}

	recipeStepID := identifiers.New()
	exampleRecipe := fakes.BuildFakeRecipe()

	exampleRecipeMedia := fakes.BuildFakeRecipeMedia()
	exampleRecipeMedia.BelongsToRecipe = &exampleRecipe.ID
	exampleRecipe.Media = []*types.RecipeMedia{
		exampleRecipeMedia,
	}

	exampleRecipePrepTask := fakes.BuildFakeRecipePrepTask()
	exampleRecipePrepTask.BelongsToRecipe = exampleRecipe.ID
	exampleRecipePrepTask.TaskSteps = []*types.RecipePrepTaskStep{
		{
			ID:                      identifiers.New(),
			BelongsToRecipeStep:     recipeStepID,
			BelongsToRecipePrepTask: exampleRecipePrepTask.ID,
			SatisfiesRecipeStep:     true,
		},
	}
	exampleRecipe.PrepTasks = []*types.RecipePrepTask{
		exampleRecipePrepTask,
	}

	exampleRecipe.Steps = []*types.RecipeStep{
		buildRecipeStepForTestCreation(t, ctx, exampleRecipe.ID, dbc),
	}
	exampleRecipePrepTask.TaskSteps[0].BelongsToRecipeStep = exampleRecipe.Steps[0].ID
	exampleRecipe.CreatedByUser = userID

	exampleRecipe.Media = []*types.RecipeMedia{}
	for i := range exampleRecipe.Steps {
		exampleRecipe.Steps[i].Media = nil                // []*types.RecipeMedia{}
		exampleRecipe.Steps[i].CompletionConditions = nil // []*types.RecipeStepCompletionCondition{}
	}

	return exampleRecipe
}

func createRecipeForTest(t *testing.T, ctx context.Context, exampleRecipe *types.Recipe, dbc *Querier, alsoCreateMeal bool) *types.Recipe {
	t.Helper()

	// create
	if exampleRecipe == nil {
		exampleRecipe = buildRecipeForTestCreation(t, ctx, "", dbc)
	}
	dbInput := converters.ConvertRecipeToRecipeDatabaseCreationInput(exampleRecipe)
	dbInput.AlsoCreateMeal = alsoCreateMeal

	created, err := dbc.CreateRecipe(ctx, dbInput)
	require.NoError(t, err)
	require.NotNil(t, created)

	recipe, err := dbc.GetRecipe(ctx, created.ID)
	require.NoError(t, err)
	require.NotNil(t, recipe)
	exampleRecipe.CreatedAt = recipe.CreatedAt

	for i := range recipe.PrepTasks {
		exampleRecipe.PrepTasks[i].CreatedAt = recipe.PrepTasks[i].CreatedAt
	}

	for i := range recipe.Steps {
		assert.Equal(t, exampleRecipe.Steps[i].Preparation.ID, recipe.Steps[i].Preparation.ID)
		exampleRecipe.Steps[i].Preparation = recipe.Steps[i].Preparation
		exampleRecipe.Steps[i].CreatedAt = recipe.Steps[i].CreatedAt

		for j := range recipe.Steps[i].Ingredients {
			exampleRecipe.Steps[i].Ingredients[j].CreatedAt = recipe.Steps[i].Ingredients[j].CreatedAt

			assert.Equal(t, exampleRecipe.Steps[i].Ingredients[j].Ingredient.ID, recipe.Steps[i].Ingredients[j].Ingredient.ID)
			exampleRecipe.Steps[i].Ingredients[j].Ingredient = recipe.Steps[i].Ingredients[j].Ingredient

			assert.Equal(t, exampleRecipe.Steps[i].Ingredients[j].MeasurementUnit.ID, recipe.Steps[i].Ingredients[j].MeasurementUnit.ID)
			exampleRecipe.Steps[i].Ingredients[j].MeasurementUnit = recipe.Steps[i].Ingredients[j].MeasurementUnit
		}
		for j := range recipe.Steps[i].Instruments {
			exampleRecipe.Steps[i].Instruments[j].CreatedAt = recipe.Steps[i].Instruments[j].CreatedAt

			assert.Equal(t, exampleRecipe.Steps[i].Instruments[j].Instrument.ID, recipe.Steps[i].Instruments[j].Instrument.ID)
			exampleRecipe.Steps[i].Instruments[j].Instrument = recipe.Steps[i].Instruments[j].Instrument
		}
		for j := range recipe.Steps[i].Vessels {
			exampleRecipe.Steps[i].Vessels[j].CreatedAt = recipe.Steps[i].Vessels[j].CreatedAt

			assert.Equal(t, exampleRecipe.Steps[i].Vessels[j].Vessel.ID, recipe.Steps[i].Vessels[j].Vessel.ID)
			exampleRecipe.Steps[i].Vessels[j].Vessel = recipe.Steps[i].Vessels[j].Vessel

			if recipe.Steps[i].Vessels[j].Vessel.CapacityUnit != nil {
				assert.Equal(t, exampleRecipe.Steps[i].Vessels[j].Vessel.CapacityUnit.ID, recipe.Steps[i].Vessels[j].Vessel.CapacityUnit.ID)
				exampleRecipe.Steps[i].Vessels[j].Vessel.CapacityUnit = recipe.Steps[i].Vessels[j].Vessel.CapacityUnit
			}
		}
		for j := range recipe.Steps[i].CompletionConditions {
			exampleRecipe.Steps[i].CompletionConditions[j].CreatedAt = recipe.Steps[i].CompletionConditions[j].CreatedAt

			assert.Equal(t, exampleRecipe.Steps[i].CompletionConditions[j].IngredientState.ID, recipe.Steps[i].CompletionConditions[j].IngredientState.ID)
			exampleRecipe.Steps[i].CompletionConditions[j].IngredientState = recipe.Steps[i].CompletionConditions[j].IngredientState

			for k := range recipe.Steps[i].CompletionConditions[j].Ingredients {
				assert.Equal(t, exampleRecipe.Steps[i].CompletionConditions[j].Ingredients[k].ID, recipe.Steps[i].CompletionConditions[j].Ingredients[k].ID)
				exampleRecipe.Steps[i].CompletionConditions[j].Ingredients[k] = recipe.Steps[i].CompletionConditions[j].Ingredients[k]
			}
		}
		for j := range recipe.Steps[i].Products {
			exampleRecipe.Steps[i].Products[j].CreatedAt = recipe.Steps[i].Products[j].CreatedAt
			assert.Equal(t, exampleRecipe.Steps[i].Products[j].MeasurementUnit.ID, recipe.Steps[i].Products[j].MeasurementUnit.ID)
			exampleRecipe.Steps[i].Products[j].MeasurementUnit = recipe.Steps[i].Products[j].MeasurementUnit
		}
		for j := range recipe.Steps[i].Media {
			assert.Equal(t, exampleRecipe.Steps[i].Media[j].ID, recipe.Steps[i].Media[j].ID)
			exampleRecipe.Steps[i].Media[j] = recipe.Steps[i].Media[j]
		}

		require.Equal(t, len(exampleRecipe.Steps[i].Products), len(recipe.Steps[i].Products))
		require.Equal(t, len(exampleRecipe.Steps[i].Instruments), len(recipe.Steps[i].Instruments))
		require.Equal(t, len(exampleRecipe.Steps[i].Vessels), len(recipe.Steps[i].Vessels))
		require.Equal(t, len(exampleRecipe.Steps[i].Ingredients), len(recipe.Steps[i].Ingredients))
		require.Equal(t, len(exampleRecipe.Steps[i].Media), len(recipe.Steps[i].Media))
		require.Equal(t, len(exampleRecipe.Steps[i].CompletionConditions), len(recipe.Steps[i].CompletionConditions))

		require.Equal(t, exampleRecipe.Steps[i], recipe.Steps[i])
	}

	require.Equal(t, exampleRecipe, recipe)

	return recipe
}

func TestQuerier_Integration_Recipes(t *testing.T) {
	if !runningContainerTests {
		t.SkipNow()
	}

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

	createdRecipes := []*types.Recipe{}

	// create
	exampleRecipe := buildRecipeForTestCreation(t, ctx, user.ID, dbc)
	createdRecipes = append(createdRecipes, createRecipeForTest(t, ctx, exampleRecipe, dbc, true))

	// update
	updatedRecipe := buildRecipeForTestCreation(t, ctx, user.ID, dbc)
	updatedRecipe.ID = createdRecipes[0].ID
	assert.NoError(t, dbc.UpdateRecipe(ctx, updatedRecipe))

	// create more
	for i := 0; i < exampleQuantity; i++ {
		exampleRecipe = buildRecipeForTestCreation(t, ctx, user.ID, dbc)
		exampleRecipe.Name = fmt.Sprintf("%s %d", exampleRecipe.Name, i)
		createdRecipes = append(createdRecipes, createRecipeForTest(t, ctx, exampleRecipe, dbc, true))
	}

	// fetch as list
	recipes, err := dbc.GetRecipes(ctx, nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, recipes.Data)
	assert.Equal(t, len(createdRecipes), len(recipes.Data))

	needingIndexing, err := dbc.GetRecipeIDsThatNeedSearchIndexing(ctx)
	assert.NoError(t, err)
	assert.NotEmpty(t, needingIndexing)

	// search
	searchResults, err := dbc.SearchForRecipes(ctx, createdRecipes[0].Name, nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, searchResults)

	byIDs, err := dbc.GetRecipesWithIDs(ctx, []string{createdRecipes[0].ID})
	assert.NoError(t, err)
	assert.NotEmpty(t, byIDs)
	assert.Equal(t, createdRecipes[0].ID, byIDs[0].ID)

	// delete
	for _, recipe := range createdRecipes {
		assert.NoError(t, dbc.MarkRecipeAsIndexed(ctx, recipe.ID))
		assert.NoError(t, dbc.ArchiveRecipe(ctx, recipe.ID, user.ID))

		var exists bool
		exists, err = dbc.RecipeExists(ctx, recipe.ID)
		assert.NoError(t, err)
		assert.False(t, exists)

		var y *types.Recipe
		y, err = dbc.GetRecipeByIDAndUser(ctx, recipe.ID, recipe.CreatedByUser)
		assert.Nil(t, y)
		assert.Error(t, err)
		assert.ErrorIs(t, err, sql.ErrNoRows)
	}
}
