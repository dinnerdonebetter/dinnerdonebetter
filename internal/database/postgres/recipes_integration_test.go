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
	preparation := createValidPreparationForTest(t, ctx, nil, dbc)
	ingredient := createValidIngredientForTest(t, ctx, nil, dbc)
	measurementUnit1 := createValidMeasurementUnitForTest(t, ctx, nil, dbc)
	recipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()
	recipeStepIngredient.Ingredient = ingredient
	recipeStepIngredient.MeasurementUnit = *measurementUnit1
	recipeStepIngredient.BelongsToRecipeStep = recipeStepID

	instrument := createValidInstrumentForTest(t, ctx, nil, dbc)
	recipeStepInstrument := fakes.BuildFakeRecipeStepInstrument()
	recipeStepInstrument.Instrument = instrument
	recipeStepInstrument.BelongsToRecipeStep = recipeStepID

	measurementUnit2 := createValidMeasurementUnitForTest(t, ctx, nil, dbc)
	exampleVessel := fakes.BuildFakeValidVessel()
	exampleVessel.CapacityUnit = measurementUnit2
	vessel := createValidVesselForTest(t, ctx, exampleVessel, dbc)
	recipeStepVessel := fakes.BuildFakeRecipeStepVessel()
	recipeStepVessel.Vessel = vessel
	recipeStepVessel.BelongsToRecipeStep = recipeStepID

	exampleRecipe.Steps = []*types.RecipeStep{
		{
			BelongsToRecipe: exampleRecipe.ID,
			ID:              recipeStepID,
			Preparation:     *preparation,
			Index:           0,
			Ingredients: []*types.RecipeStepIngredient{
				recipeStepIngredient,
			},
			Instruments: []*types.RecipeStepInstrument{
				recipeStepInstrument,
			},
			Vessels: []*types.RecipeStepVessel{
				recipeStepVessel,
			},
		},
	}
	exampleRecipe.Media = []*types.RecipeMedia{}
	exampleRecipe.PrepTasks = []*types.RecipePrepTask{}
	exampleRecipe.CreatedByUser = userID

	return exampleRecipe
}

func createRecipeForTest(t *testing.T, ctx context.Context, exampleRecipe *types.Recipe, dbc *Querier) *types.Recipe {
	t.Helper()

	// create
	if exampleRecipe == nil {
		exampleRecipe = buildRecipeForTestCreation(t, ctx, "", dbc)
	}
	dbInput := converters.ConvertRecipeToRecipeDatabaseCreationInput(exampleRecipe)
	dbInput.AlsoCreateMeal = true

	created, err := dbc.CreateRecipe(ctx, dbInput)
	require.NoError(t, err)
	require.NotNil(t, created)

	recipe, err := dbc.GetRecipe(ctx, created.ID)
	require.NoError(t, err)
	require.NotNil(t, recipe)

	exampleRecipe.CreatedAt = recipe.CreatedAt
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
			exampleRecipe.Steps[i].Media[j].CreatedAt = recipe.Steps[i].Media[j].CreatedAt
		}
	}

	assert.Equal(t, exampleRecipe, recipe)

	return recipe
}

func TestQuerier_Integration_Recipes(t *testing.T) {
	if !runningContainerTests {
		t.SkipNow()
	}

	ctx := context.Background()
	dbc, container := buildDatabaseClientForTest(t, ctx)

	defer func(t *testing.T) {
		t.Helper()
		assert.NoError(t, container.Terminate(ctx))
	}(t)

	user := createUserForTest(t, ctx, nil, dbc)

	createdRecipes := []*types.Recipe{}

	// create
	exampleRecipe := buildRecipeForTestCreation(t, ctx, user.ID, dbc)
	createdRecipes = append(createdRecipes, createRecipeForTest(t, ctx, exampleRecipe, dbc))

	// update
	updatedRecipe := buildRecipeForTestCreation(t, ctx, user.ID, dbc)
	updatedRecipe.ID = createdRecipes[0].ID
	assert.NoError(t, dbc.UpdateRecipe(ctx, updatedRecipe))

	// create more
	for i := 0; i < exampleQuantity; i++ {
		exampleRecipe = buildRecipeForTestCreation(t, ctx, user.ID, dbc)
		exampleRecipe.Name = fmt.Sprintf("%s %d", exampleRecipe.Name, i)
		createdRecipes = append(createdRecipes, createRecipeForTest(t, ctx, exampleRecipe, dbc))
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

	// delete
	for _, recipe := range createdRecipes {
		assert.NoError(t, dbc.MarkRecipeAsIndexed(ctx, recipe.ID))
		assert.NoError(t, dbc.ArchiveRecipe(ctx, recipe.ID, user.ID))

		var exists bool
		exists, err = dbc.RecipeExists(ctx, recipe.ID)
		assert.NoError(t, err)
		assert.False(t, exists)

		var y *types.Recipe
		y, err = dbc.GetRecipe(ctx, recipe.ID)
		assert.Nil(t, y)
		assert.Error(t, err)
		assert.ErrorIs(t, err, sql.ErrNoRows)
	}
}
